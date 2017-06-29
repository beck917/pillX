package pillx

import (
	"fmt"
	"hash/adler32"
	"strconv"
	"sync"
	"time"

	"github.com/beck917/pillX/Proto"

	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/clientv3"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

type Gateway struct {
	InnerAddr     string
	OuterAddr     string
	OuterProtocol IProtocol
	InnerServer   *Server
	OuterServer   *Server
	EtcdClient    *etcd.Client
	GatewayName   string
	WatchName     string
}

//var workers map[uint32]*Response
var worker_id uint32 = 0

var (
	clientsMu           sync.RWMutex
	secretKey                                = ")*^BHFifd2/*-fdsBHJ98"
	workerPools         map[string]Pool      = make(map[string]Pool)
	clients             map[uint64]*Response = make(map[uint64]*Response)
	chat_channel        *Channel             = NewChannel("chat")
	channels            map[int32]*Channel   = make(map[int32]*Channel)
	workerHandshakeFlgs                      = make(map[uint64]bool)
	uidBindClientId                          = make(map[int32]uint64) //uid和clientid关联
	clientMap                                = make(map[uint64]*ClientData)
	//禁言
	banids = make(map[int32]*userBaned)
	//踢人
	kickids = make(map[int32]*userBaned)
	//黑名单
	BlackIdMap       = make(map[int32]bool)
	AdminIdMap       = make(map[int32]bool)
	bantime    int64 = 300  //禁言5分钟
	kicktime   int64 = 1800 //踢出30分钟不准进入

	//最近聊天内容纪录
	RecentMsgChan = NewMsgList(256)

	//广播某某进入，某某退出

	//chatdata 加入时间
	//日志系统完善
)

type userBaned struct {
	uid       int32
	timestamp int64 //屏蔽时间
}

type ClientData struct {
	Uid      int32
	RoomId   int32
	Channel  *Channel
	Platform int32
	Ip       string
	Uname    string
}

func NewGateway(innerAddr string, outProtocol IProtocol) *Gateway {
	gateway := &Gateway{
		InnerAddr:     innerAddr,
		OuterProtocol: outProtocol,
	}
	gateway.Init()
	return &Gateway{}
}

func (gateway *Gateway) innerConnectHandler(worker *Response, protocol IProtocol) {
	//worker_id = atomic.AddUint32(&worker_id, 1)
	//workers[worker_id] = worker
	//fmt.Printf("worker %d 连接到此网关\n", worker_id)
}

func (gateway *Gateway) innerMessageHandler(worker *Response, protocol IProtocol) {
	//将gateway协议转化为客户端协议
	pillProtocol := &PillProtocol{}
	header := &PillProtocolHeader{}
	pillProtocol.Header = header

	req := protocol.(*GateWayProtocol)
	header.Cmd = req.Header.Cmd
	header.Error = req.Header.Error
	header.Mark = req.Header.Mark
	header.Size = req.Header.Size
	header.Version = PILL_VERSION
	header.Sid = req.Header.Sid
	pillProtocol.Content = req.Content

	//发送给client
	//TODO 这里用了辣鸡锁，后期改成用nil判断
	clientsMu.Lock()
	client, ok := clients[req.Header.ClientId]
	clientsMu.Unlock()
	if !ok {
		//如果没有此key，说明被删除了
		return
	}
	ret := gateway.returnMsg(req.Header.Cmd, req.Header.Sid, 0x11, "ok", 0)
	client.Send(ret)

	//记录消息
	chatData := &Proto.ChatData{}
	proto.Unmarshal(pillProtocol.Content, chatData)
	RecentMsgChan.Push(chatData)

	//广播消息
	chat_channel.Publish(client, pillProtocol)

	MyLog().WithFields(log.Fields{
		"client_id": req.Header.ClientId,
		"content":   string(pillProtocol.Content),
		"client_ip": client.GetConn().remonte_conn.RemoteAddr(),
		"room_id":   clientMap[req.Header.ClientId].RoomId,
		"platform":  clientMap[req.Header.ClientId].Platform,
	}).Info("广播完毕")
}

func (gateway *Gateway) innerCloseHandler(worker *Response, protocol IProtocol) {

}

func (gateway *Gateway) outerConnectHandler(client *Response, protocol IProtocol) {
	clients[client.GetConn().Id] = client
	MyLog().WithFields(log.Fields{
		"client_id": client.GetConn().Id,
		"client_ip": client.GetConn().remonte_conn.RemoteAddr(),
	}).Info("连接到网关")
}

func (gateway *Gateway) outerHandShakeHandler(client *Response, protocol IProtocol) {
	req := protocol.(*PillProtocol)
	MyLog().WithFields(log.Fields{
		"client_id": client.GetConn().Id,
		"info_code": "hs1",
	}).Info("握手协议开始")

	//被踢掉用户
	if _, okid := clientMap[client.GetConn().Id]; okid {
		if userban, ok := kickids[clientMap[client.GetConn().Id].Uid]; ok {
			if userban.timestamp+kicktime > time.Now().Unix() {
				ret := gateway.returnMsg(req.Header.Cmd, req.Header.Sid, 0xf2, "你被管理员点名，30分钟内无法进入！", 0xf2)
				client.Send(ret)
				client.Close()
				return
			} else {
				delete(banids, clientMap[client.GetConn().Id].Uid)
			}
		}
	}

	/**
	//取出timestamp
	timestamp := binary.BigEndian.Uint32(req.Content[:4])
	//取出token
	token := binary.BigEndian.Uint32(req.Content[4:])
	*/

	//数据处理到proto类
	handShakeData := &Proto.HandShake{}
	proto.Unmarshal(req.Content, handShakeData)

	timestamp := handShakeData.GetTimestamp()
	token := handShakeData.GetToken()

	//验证token
	secretStr := fmt.Sprintf("%d%d%d%s", req.Header.Cmd, timestamp, req.Header.Sid, secretKey)
	MyLog().Info(req.Header.Cmd, timestamp, req.Header.Sid, secretKey)
	MyLog().Info(secretStr)
	serverToken := adler32.Checksum([]byte(secretStr))
	MyLog().Info(serverToken)
	MyLog().Info(token)
	if serverToken == uint32(token) {
		MyLog().WithFields(log.Fields{
			"client_id": client.GetConn().Id,
			"msgcode":   "hs2",
		}).Info("握手验证成功")
		//验证通过，修改conn中的
		client.GetConn().HandshakeFlg = true
	} else {
		//req.Header.Error = SYS_CONNECT_HANDSHAKE_ERROR
		MyLog().WithFields(log.Fields{
			"client_id": client.GetConn().Id,
			"info_code": "hs2",
		}).Info("握手验证失败")
		//返回错误
		reterr := gateway.returnMsg(req.Header.Cmd, req.Header.Sid, 0xf1, "握手验证失败", SYS_CONNECT_HANDSHAKE_ERROR)
		client.Send(reterr)

		//断开
		client.Close()
		return
	}

	//创建频道
	if channels[handShakeData.GetRoomId()] == nil {
		channels[handShakeData.GetRoomId()] = NewChannel("chat" + strconv.Itoa(int(handShakeData.GetRoomId())))
	}
	//订阅
	channels[handShakeData.GetRoomId()].Subscribe(client)

	//订阅全部频道
	chat_channel.Subscribe(client)

	//推送到此频道的所有人通知

	//绑定uid和clientid的关系
	if handShakeData.Uid != nil {
		uidBindClientId[handShakeData.GetUid()] = client.conn.Id
	}

	//创建client数据
	clientData := &ClientData{
		Uid:      handShakeData.GetUid(),
		RoomId:   handShakeData.GetRoomId(),
		Channel:  channels[handShakeData.GetRoomId()],
		Platform: handShakeData.GetPlatform(),
		Ip:       client.conn.remote_addr,
		Uname:    handShakeData.GetUname(),
	}
	clientMap[client.conn.Id] = clientData

	//发送握手到worker
	//worker := workerPool["123"].Get()
	//worker.Send(req)

	//广播xxx进入
	/**
	pillin := NewPillProtocol()
	pillin.Header.Cmd = SYS_ON_CLIENTIN
	chatProto := &Proto.ChatData{}
	chatProto.Msg = fmt.Sprintf("%s进入情报大厅", string(clientData.Uname))
	pillin.Content, _ = proto.Marshal(chatProto)
	pillin.Header.Size = uint16(len(pillin.Content))
	chat_channel.Publish(client, pillin)


	MyLog().WithFields(log.Fields{
		"client_id": client.GetConn().Id,
		"platform":  handShakeData.GetPlatform(),
		"header":    pillin.Header,
		"content":   string(pillin.Content),
		"uname":     clientData.Uname,
	}).Info("广播完毕")
	*/
	ret := gateway.returnMsg(req.Header.Cmd, req.Header.Sid, 0x11, "握手成功", 0)
	client.Send(ret)

	MyLog().WithFields(log.Fields{
		"client_id": client.GetConn().Id,
		"platform":  handShakeData.GetPlatform(),
		"header":    ret.Header,
		"content":   string(ret.Content),
	}).Info("握手返回")

	MyLog().WithFields(log.Fields{
		"client_id":     client.GetConn().Id,
		"connnet_count": len(clientMap),
	}).Info("统计信息")
}

func (gateway *Gateway) outerMessageHandler(client *Response, protocol IProtocol) {
	//握手，block请求忽略
	req := protocol.(*PillProtocol)
	if req.GetCmd() == SYS_ON_HANDSHAKE || req.GetCmd() == SYS_ON_BLOCK || req.GetCmd() == SYS_ON_HEARTBEAT || req.GetCmd() == SYS_ON_KICK || req.GetCmd() == SYS_ON_BLACK {
		return
	}

	if client.GetConn().HandshakeFlg == false {
		//验证是否握手过
		ret := gateway.returnMsg(req.Header.Cmd, req.Header.Sid, 0xf3, "没有握手", 0xf3)
		client.Send(ret)
		client.Close()
	}

	//ban禁言和black用户
	if _, cok := clientMap[client.GetConn().Id]; !cok {
		ret := gateway.returnMsg(req.Header.Cmd, req.Header.Sid, 0xf4, "用户信息不完整", 0xf4)
		client.Send(ret)
	}

	if userban, ok := banids[clientMap[client.GetConn().Id].Uid]; ok {
		if userban.timestamp+bantime > time.Now().Unix() {
			ret := gateway.returnMsg(req.Header.Cmd, req.Header.Sid, 0xf2, "禁止发言", 0xf2)
			client.Send(ret)
			return
		} else {
			delete(banids, clientMap[client.GetConn().Id].Uid)
		}
	}

	if BlackIdMap[clientMap[client.GetConn().Id].Uid] == true {
		ret := gateway.returnMsg(req.Header.Cmd, req.Header.Sid, 0xf2, "已经进入小黑屋", 0xf2)
		client.Send(ret)
	}

	//将客户端协议转化为gateway协议
	gatewayProtocol := &GateWayProtocol{}
	header := &GatewayHeader{}
	gatewayProtocol.Header = header

	header.ClientId = client.GetConn().Id

	header.Cmd = req.Header.Cmd
	header.Error = req.Header.Error
	header.Mark = PROTO_HEADER_FIRSTCHAR
	header.Version = GATEWAY_VERSION
	header.Sid = req.Header.Sid
	header.Size = uint16(len(req.Content))
	gatewayProtocol.Content = req.Content

	//发送给一个合适的worker,根据clientid做hash
	workerPool, workerKey := GetPool(workerPools)
	if workerPool == nil {
		return
	}
	/**
	if workerPool == nil {
		log.Error("worker池未初始化")
		errorMsg := &PillProtocolHeader{
			Mark:  PROTO_HEADER_FIRSTCHAR,
			Size:  0,
			Cmd:   0x0001,
			Error: 1,
		}
	}
	*/
	worker, err := workerPool.Get()
	if err != nil {
		MyLog().WithError(err).Error("worker池返回错误")
		return
	}
	//发送握手
	if !workerHandshakeFlgs[header.ClientId] {
		handshakeProto := &Proto.WorkerHandShark{}
		handshakeProto.IP = proto.String(client.conn.remote_addr)
		handshakeGateway := NewGatewayProtocol()
		handshakeGateway.Content, _ = proto.Marshal(handshakeProto)
		handshakeGateway.Header.Size = uint16(len(handshakeGateway.Content))
		handshakeGateway.Header.Cmd = SYS_ON_HANDSHAKE
		handshakeGateway.Header.ClientId = header.ClientId
		MyLog().WithField("proto", handshakeGateway.Header).Info("发送握手信息给worker ", workerKey)
		worker.response.Send(handshakeGateway)

		workerHandshakeFlgs[header.ClientId] = true
	}

	_, err = worker.response.Send(gatewayProtocol)
	if err != nil {
		//连接写入出错，记录错误信息
		MyLog().WithField("proto", gatewayProtocol.Header).Error(err)
		pillerror := NewPillProtocol()
		pillerror.Header.Error = SYS_CONNECT_WORKER_ERROR
		client.Send(pillerror)
	}

	//回收连接
	worker.Close()
	//fmt.Printf("%x", gatewayProtocol.Header)
	//fmt.Println(gatewayProtocol.Header)
	//fmt.Printf("%x", req.Header)
	//fmt.Printf("%s", req.Content)

	MyLog().WithFields(log.Fields{
		"client_id": client.GetConn().Id,
		"content":   string(gatewayProtocol.Content),
		"header":    gatewayProtocol.Header,
		"client_ip": client.GetConn().remonte_conn.RemoteAddr(),
		"room_id":   clientMap[client.GetConn().Id].RoomId,
		"platform":  clientMap[client.GetConn().Id].Platform,
	}).Info("发送给worker ", workerKey)
}

func (gateway *Gateway) outerBlockHandle(client *Response, protocol IProtocol) {
	//取出频道
	channel := clientMap[client.conn.Id].Channel

	req := protocol.(*PillProtocol)

	//解析content
	//数据处理到proto类
	blockData := &Proto.BlockData{}
	proto.Unmarshal(req.Content, blockData)

	/**
	MyLog().Info(blockData.GetBlockUid(), uidBindClientId[blockData.GetBlockUid()])

	if _, ok := uidBindClientId[blockData.GetBlockUid()]; ok {
		channel.Block(client, uidBindClientId[blockData.GetBlockUid()])
		chat_channel.Block(client, uidBindClientId[blockData.GetBlockUid()])
	}
	*/

	channel.BlockUid(clientMap[client.conn.Id].Uid, blockData.GetBlockUid())
	chat_channel.BlockUid(clientMap[client.conn.Id].Uid, blockData.GetBlockUid())

	//返回成功消息
	ret := gateway.returnMsg(req.Header.Cmd, req.Header.Sid, 0x11, "屏蔽成功", 0)
	client.Send(ret)
}

//禁言
func (gateway *Gateway) outerBanHandle(client *Response, protocol IProtocol) {
	req := protocol.(*PillProtocol)

	//数据处理到proto类
	blockData := &Proto.BlockData{}
	proto.Unmarshal(req.Content, blockData)
	uid := blockData.GetBlockUid()

	userban := &userBaned{
		uid:       uid,
		timestamp: time.Now().Unix(),
	}
	banids[uid] = userban

	//返回成功消息
	ret := gateway.returnMsg(req.Header.Cmd, req.Header.Sid, 0x11, "禁言成功", 0)
	client.Send(ret)
}

//踢人
func (gateway *Gateway) outerKickHandle(client *Response, protocol IProtocol) {
	req := protocol.(*PillProtocol)

	//数据处理到proto类
	blockData := &Proto.BlockData{}
	proto.Unmarshal(req.Content, blockData)
	uid := blockData.GetBlockUid()

	userban := &userBaned{
		uid:       uid,
		timestamp: time.Now().Unix(),
	}
	kickids[uid] = userban

	//断开链接
	clients[uidBindClientId[uid]].Close()

	//返回
	ret := gateway.returnMsg(req.Header.Cmd, req.Header.Sid, 0x11, "踢掉成功", 0)
	client.Send(ret)
}

//黑名单
func (gateway *Gateway) outerBLackHandle(client *Response, protocol IProtocol) {
	req := protocol.(*PillProtocol)

	//数据处理到proto类
	blockData := &Proto.BlockData{}
	proto.Unmarshal(req.Content, blockData)
	uid := blockData.GetBlockUid()

	BlackIdMap[uid] = true

	//返回成功消息
	ret := gateway.returnMsg(req.Header.Cmd, req.Header.Sid, 0x11, "加入黑名单成功", 0)
	client.Send(ret)
}

func (gateway *Gateway) outerHeartBeatHandle(client *Response, protocol IProtocol) {
	//返回成功消息
	req := protocol.(*PillProtocol)
	req.Content = []byte("pong")
	req.Header.Size = 4

	client.Send(req)
}

func (gateway *Gateway) outerCloseHandler(client *Response, protocol IProtocol) {

	//日志
	MyLog().WithFields(log.Fields{
		"client_id": client.GetConn().Id,
		"client_ip": client.GetConn().remonte_conn.RemoteAddr(),
		"info_code": "close1",
	}).Info("连接断开")

	//各种清除
	clientsMu.Lock()
	delete(clients, client.conn.Id)
	if _, ok := clientMap[client.conn.Id]; ok {
		delete(uidBindClientId, clientMap[client.conn.Id].Uid)
	}
	delete(clientMap, client.conn.Id)
	clientsMu.Unlock()
}

func (gateway *Gateway) returnMsg(cmd uint16, sid uint32, code int32, msg string, errcode uint16) (pillret *PillProtocol) {
	pillret = NewPillProtocol()
	pillret.Header.Cmd = cmd
	messageProto := &Proto.MessageData{}
	messageHeaderProto := &Proto.ChatHeader{}
	messageHeaderProto.Code = proto.Int32(code)
	messageHeaderProto.Msg = proto.String(msg)
	messageProto.Header = messageHeaderProto
	pillret.Content, _ = proto.Marshal(messageProto)
	pillret.Header.Size = uint16(len(pillret.Content))
	pillret.Header.Error = errcode
	pillret.Header.Sid = sid

	//记录返回日志
	MyLog().WithFields(log.Fields{
		"header": pillret.Header,
		"code":   code,
		"msg":    msg,
	}).Info("返回记录")
	return
}

func (gateway *Gateway) Init() {
	gateway.InnerServer = &Server{
		Addr:     gateway.InnerAddr,
		Protocol: &GateWayProtocol{},
	}
	gateway.InnerServer.HandleFunc(SYS_ON_CONNECT, gateway.innerConnectHandler)
	gateway.InnerServer.HandleFunc(SYS_ON_MESSAGE, gateway.innerMessageHandler)
	gateway.InnerServer.HandleFunc(SYS_ON_CLOSE, gateway.innerCloseHandler)

	go gateway.InnerServer.ListenAndServe()
	MyLog().WithFields(log.Fields{
		"addr": gateway.InnerAddr,
	}).Info("inner server started")

	gateway.OuterServer = &Server{
		Addr:     gateway.OuterAddr,
		Handler:  nil,
		Protocol: gateway.OuterProtocol,
	}
	gateway.OuterServer.HandleFunc(SYS_ON_CONNECT, gateway.outerConnectHandler)
	gateway.OuterServer.HandleFunc(SYS_ON_MESSAGE, gateway.outerMessageHandler)
	gateway.OuterServer.HandleFunc(SYS_ON_CLOSE, gateway.outerCloseHandler)
	gateway.OuterServer.HandleFunc(SYS_ON_HANDSHAKE, gateway.outerHandShakeHandler)
	gateway.OuterServer.HandleFunc(SYS_ON_BLOCK, gateway.outerBlockHandle)
	gateway.OuterServer.HandleFunc(SYS_ON_KICK, gateway.outerKickHandle)
	gateway.OuterServer.HandleFunc(SYS_ON_BLACK, gateway.outerBLackHandle)
	gateway.OuterServer.HandleFunc(SYS_ON_BAN, gateway.outerBanHandle)
	gateway.OuterServer.HandleFunc(SYS_ON_HEARTBEAT, gateway.outerHeartBeatHandle)
	go gateway.OuterServer.ListenAndServe()
	MyLog().WithFields(log.Fields{
		"addr": gateway.OuterAddr,
	}).Info("outer server started")
	//注册gateway
	gateway.EtcdClient.Put(context.Background(), gateway.GatewayName, gateway.InnerAddr)

	//监听etcd worker注册
	MyLog().Info("etcd watch started")
	rch := gateway.EtcdClient.Watch(context.Background(), gateway.WatchName, etcd.WithPrefix())
	for wresp := range rch {
		if wresp.Events != nil {
			gateway.watchWorkers(wresp.Events)
		}
	}
	MyLog().Info("etcd watch started2")
	return
}

func (gateway *Gateway) watchWorkers(events []*etcd.Event) {
	for _, ev := range events {
		name := string(ev.Kv.Key)
		value := string(ev.Kv.Value)
		MyLog().WithFields(log.Fields{
			"type":  ev.Type,
			"key":   string(ev.Kv.Key),
			"value": string(ev.Kv.Value),
		}).Info("workers信息")
		//MyLog().Info(ev.Type)
		if ev.Type == 1 {
			//清除此节点
			//workerPools[name] = nil
			delete(workerPools, name)
			MyLog().Info("worker deleted ", name)
		} else {
			if _, ok := workerPools[name]; !ok {
				client := &Client{}
				client.Addr = value
				workerPools[name], _ = client.Dail()
				MyLog().Info("worker connected ", name)
			}
		}
		//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
	}
	return
}

func (gateway *Gateway) ConnectWorkers() {
	return
}
