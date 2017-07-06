package pillx

import (
	"strconv"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/beck917/pillX/Proto"
	etcd "github.com/coreos/etcd/clientv3"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

type GatewayWebsocket struct {
	InnerAddr     string
	OuterAddr     string
	OuterProtocol IProtocol
	InnerServer   *Server
	OuterServer   *Server
	EtcdClient    *etcd.Client
	GatewayName   string
	WatchName     string
}

var (
	workersMu     sync.RWMutex
	workers       map[string]*Response = make(map[string]*Response)
	clientSyncMap                      = struct {
		sync.RWMutex
		data map[uint64]*ClientDataStruct
	}{data: make(map[uint64]*ClientDataStruct)}
	workerHandshakeMap = struct {
		sync.RWMutex
		data map[uint64]string
	}{data: make(map[uint64]string)}
)

type ClientDataStruct struct {
	Uid      int
	RoomId   int
	Channel  *Channel
	Platform int
	Ip       string
	Uname    string
}

func (gateway *GatewayWebsocket) innerConnectHandler(worker *Response, protocol IProtocol) {
	//worker_id = atomic.AddUint32(&worker_id, 1)
	//workers[worker_id] = worker
	//fmt.Printf("worker %d 连接到此网关\n", worker_id)
	//存储此worker的id
	workersMu.Lock()
	defer workersMu.Unlock()
	workers[worker.conn.remonte_conn.RemoteAddr().String()] = worker
	MyLog().Info("connect", worker.conn.remonte_conn.RemoteAddr().String())
}

func (gateway *GatewayWebsocket) innerCloseHandler(worker *Response, protocol IProtocol) {
	workersMu.Lock()
	defer workersMu.Unlock()
	remoteAddr := worker.conn.remonte_conn.RemoteAddr().String()
	//删除一致性哈希
	globalConsistent.Remove(remoteAddr)
	delete(workers, remoteAddr)

	//重新进行握手
	gateway.reHashConnection()
	MyLog().Info("disconnect", worker.conn.remonte_conn.RemoteAddr().String())
}

func (gateway *GatewayWebsocket) innerMessageHandler(worker *Response, protocol IProtocol) {
	//将gateway协议转化为客户端协议
	req := protocol.(*GateWayProtocol)

	//判断协议
	var outProtocol IProtocol
	switch gateway.OuterProtocol.(type) {
	case *WebSocketProtocol:
		outProtocol = &WebSocketProtocol{}
		outProtocol.(*WebSocketProtocol).Header = &WebSocketHeader{}
		outProtocol.(*WebSocketProtocol).Content = req.Content
	case *PillProtocol:
		outProtocol = &PillProtocol{}
		outProtocol.(*PillProtocol).Header = &PillProtocolHeader{}
		outProtocol.(*PillProtocol).Content = req.Content
	}

	//发送给client
	//TODO 这里用了辣鸡锁，后期改成用nil判断
	clientsMu.Lock()
	client, ok := clients[req.Header.ClientId]
	clientsMu.Unlock()
	if !ok {
		//如果没有此key，说明被删除了
		return
	}

	switch req.Header.Cmd {
	case SYS_WORKER_DISCONNECT:
		//delete(workerPools, worker)
	case SYS_BIND_UID:
		//将client_id与uid绑定
		//clientMap[req.Header.ClientId].
		clientSyncMap.Lock()
		clientSyncMap.data[req.Header.ClientId].Uid, _ = strconv.Atoi(string(req.Content))
		clientSyncMap.Unlock()
	case SYS_SEND_TO_ONE:
		client.Send(outProtocol)
	default:
		client.Send(outProtocol)
	}
	/**
	    switch cmd {
	        case GatewayProtocol::CMD_WORKER_CONNECT:
	        case GatewayProtocol::CMD_SEND_TO_ONE:
	        case GatewayProtocol::CMD_KICK:
	        case GatewayProtocol::CMD_SEND_TO_ALL:
	            // 更新客户端session
	        case GatewayProtocol::CMD_UPDATE_SESSION:
	            // 获得客户端在线状态 Gateway::getALLClientInfo()
	        case GatewayProtocol::CMD_GET_ALL_CLIENT_INFO:
	            // 判断某个client_id是否在线 Gateway::isOnline($client_id)
	        case GatewayProtocol::CMD_IS_ONLINE:
	            // 将client_id与uid绑定
	        case GatewayProtocol::CMD_BIND_UID:
	            // client_id与uid解绑 Gateway::unbindUid($client_id, $uid);
	        case GatewayProtocol::CMD_UNBIND_UID:
	            // 发送数据给uid Gateway::sendToUid($uid, $msg);
	        case GatewayProtocol::CMD_SEND_TO_UID:
	            // 将$client_id加入用户组 Gateway::joinGroup($client_id, $group);
	        case GatewayProtocol::CMD_JOIN_GROUP:
	            // 将$client_id从某个用户组中移除 Gateway::leaveGroup($client_id, $group);
	        case GatewayProtocol::CMD_LEAVE_GROUP:
	            // 向某个用户组发送消息 Gateway::sendToGroup($group, $msg);
	        case GatewayProtocol::CMD_SEND_TO_GROUP:
	            // 获取某用户组成员信息 Gateway::getClientInfoByGroup($group);
	        case GatewayProtocol::CMD_GET_CLINET_INFO_BY_GROUP:
	            // 获取用户组成员数 Gateway::getClientCountByGroup($group);
	        case GatewayProtocol::CMD_GET_CLIENT_COUNT_BY_GROUP:
	            // 获取与某个uid绑定的所有client_id Gateway::getClientIdByUid($uid);
	        case GatewayProtocol::CMD_GET_CLIENT_ID_BY_UID:
		}
	*/

	//广播消息
	//chat_channel.Publish(client, wsProtocol)

	MyLog().WithFields(log.Fields{
		"client_id": req.Header.ClientId,
		"cmd":       req.Header.Cmd,
		//"content":   string(req.Content),
		"client_ip": client.GetConn().remonte_conn.RemoteAddr(),
		//"room_id":   clientMap[req.Header.ClientId].RoomId,
		//"platform":  clientMap[req.Header.ClientId].Platform,
	}).Info("广播完毕")
}

func (gateway *GatewayWebsocket) outerConnectHandler(client *Response, protocol IProtocol) {
	clients[client.GetConn().Id] = client
	MyLog().WithFields(log.Fields{
		"client_id": client.GetConn().Id,
		"client_ip": client.GetConn().remonte_conn.RemoteAddr(),
	}).Info("连接到网关")

	//订阅全部频道
	chat_channel.Subscribe(client)

	//发送握手
}

func (gateway *GatewayWebsocket) outerMessageHandler(client *Response, protocol IProtocol) {
	//判断协议
	var outProtocol IProtocol
	var outContent []byte
	switch gateway.OuterProtocol.(type) {
	case *WebSocketProtocol:
		outProtocol = protocol.(*WebSocketProtocol)
		outContent = outProtocol.(*WebSocketProtocol).Content
	case *PillProtocol:
		outProtocol = protocol.(*PillProtocol)
		outContent = outProtocol.(*PillProtocol).Content
	}
	//req := protocol.(*WebSocketProtocol)

	MyLog().WithFields(log.Fields{
		"client_id": client.GetConn().Id,
		"content":   string(outContent),
		"client_ip": client.GetConn().remonte_conn.RemoteAddr(),
	}).Info("发送给worker")

	/**
	jsonObj, jsonErr := simplejson.NewJson(outProtocol.Content)
	if jsonErr != nil {
		//记录错误
		return
	}

	method, err := jsonObj.Get("method").String()

	if err != nil {
		return
	}
	*/

	//将客户端协议转化为gateway协议
	gatewayProtocol := &GateWayProtocol{}
	header := &GatewayHeader{}
	gatewayProtocol.Header = header

	header.ClientId = client.GetConn().Id
	header.Cmd = 0 //utils.Crc16([]byte(method))
	header.Error = 0
	header.Mark = PROTO_HEADER_FIRSTCHAR
	header.Version = GATEWAY_VERSION
	header.Sid = 0
	header.Size = uint16(len(outContent))
	gatewayProtocol.Content = outContent

	/**
	//发送给一个合适的worker,根据clientid做hash
	_, worker := GetResponse(workers)
	if worker == nil {
		MyLog().Error("no available worker")
		return
	}

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
	//初始化clientmap
	clientSyncMap.Lock()
	clientSyncMap.data[header.ClientId] = &ClientDataStruct{}
	clientSyncMap.Unlock()

	//发送握手
	serverName, res := GetResponse(header.ClientId, workers)
	if workerHandshakeMap.data[header.ClientId] != serverName {
		handshakeProto := &Proto.WorkerHandShark{}
		handshakeProto.IP = proto.String(client.conn.remote_addr)
		handshakeProto.Uid = proto.Int32(int32(0))
		handshakeGateway := NewGatewayProtocol()
		handshakeGateway.Content, _ = proto.Marshal(handshakeProto)
		handshakeGateway.Header.Size = uint16(len(handshakeGateway.Content))
		handshakeGateway.Header.Cmd = SYS_ON_HANDSHAKE
		handshakeGateway.Header.ClientId = header.ClientId
		MyLog().WithField("proto", handshakeGateway.Header).Info("发送握手信息给worker ", serverName)
		res.Send(handshakeGateway)
		workerHandshakeMap.data[header.ClientId] = serverName
	}

	_, err := responseSend(header.ClientId, workers, gatewayProtocol)
	if err != nil {
		//连接写入出错，记录错误信息
		MyLog().WithField("proto", gatewayProtocol.Header).Error(err)
		//pillerror := NewPillProtocol()
		//pillerror.Header.Error = SYS_CONNECT_WORKER_ERROR
		//client.Send(pillerror)
	}

	//fmt.Printf("%x", gatewayProtocol.Header)

	MyLog().WithFields(log.Fields{
		"client_id": client.GetConn().Id,
		"content":   string(gatewayProtocol.Content),
		"header":    gatewayProtocol.Header,
		"client_ip": client.GetConn().remonte_conn.RemoteAddr(),
	}).Info("发送给worker")

	//jsonObj.Set("content", "online")
	//req.Content, _ = jsonObj.Encode()
	//广播消息
	//chat_channel.Publish(client, req)
}

func (gateway *GatewayWebsocket) outerCloseHandler(client *Response, protocol IProtocol) {
	//日志
	MyLog().WithFields(log.Fields{
		"client_id": client.GetConn().Id,
		"client_ip": client.GetConn().remonte_conn.RemoteAddr(),
		"info_code": "close1",
	}).Info("连接断开")

	//告知worker
	closeGateway := NewGatewayProtocol()
	closeGateway.Header.Cmd = SYS_CLIENT_DISCONNECT
	closeGateway.Header.ClientId = client.GetConn().Id
	_, err := responseSend(closeGateway.Header.ClientId, workers, closeGateway)
	if err != nil {
		//连接写入出错，记录错误信息
		MyLog().WithField("proto", closeGateway.Header).Error(err)
	}
	MyLog().WithField("proto", closeGateway.Header).Info("发送关闭信息给worker")

	//各种清除
	clientsMu.Lock()
	defer clientsMu.Unlock()
	delete(clients, client.conn.Id)
}

func (gateway *GatewayWebsocket) reHashConnection() {
	//遍历客户端信息
	for id, client := range clients {
		clientIdStr := strconv.Itoa(int(id))
		serverName, _ := globalConsistent.Get(clientIdStr)

		if serverName == "" {
			break
		}

		MyLog().Info("servet name ", serverName)
		MyLog().Info("whm ", workerHandshakeMap.data[id])
		MyLog().Info("clientSyncMap ", clientSyncMap.data[id].Uid)

		//判断节点是否发生了变化,如果变化了,则重新发送握手信息
		if workerHandshakeMap.data[id] != "" && workerHandshakeMap.data[id] != serverName {
			if clientSyncMap.data[id] != nil && clientSyncMap.data[id].Uid != 0 {
				gateway.sendHandShakeMsg(client, clientSyncMap.data[id].Uid)
			}
		}
	}
}

func (gateway *GatewayWebsocket) sendHandShakeMsg(client *Response, uid int) {
	handshakeProto := &Proto.WorkerHandShark{}
	handshakeProto.IP = proto.String(client.conn.remote_addr)
	handshakeProto.Uid = proto.Int32(int32(uid))
	handshakeGateway := NewGatewayProtocol()
	handshakeGateway.Content, _ = proto.Marshal(handshakeProto)
	handshakeGateway.Header.Size = uint16(len(handshakeGateway.Content))
	handshakeGateway.Header.Cmd = SYS_ON_HANDSHAKE
	handshakeGateway.Header.ClientId = client.conn.Id

	serverName, res := GetResponse(handshakeGateway.Header.ClientId, workers)
	MyLog().WithField("proto", handshakeGateway.Header).Info("发送握手信息给worker ", serverName)
	res.Send(handshakeGateway)

	workerHandshakeMap.data[handshakeGateway.Header.ClientId] = serverName
}

func (gateway *GatewayWebsocket) Init() {
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
	//gateway.OuterServer.HandleFunc(SYS_ON_HANDSHAKE, gateway.outerHandShakeHandler)
	go gateway.OuterServer.ListenAndServe()

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
}

func (gateway *GatewayWebsocket) watchWorkers(events []*etcd.Event) {
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
			//delete(workers, value)
			MyLog().Info("worker deleted ", name, value)
		} else {
			if _, ok := workers[value]; !ok {
				workersMu.Lock()
				defer workersMu.Unlock()
				client := NewGatewayClient(value)
				client.HandleFunc(SYS_ON_CONNECT, gateway.innerConnectHandler)
				client.HandleFunc(SYS_ON_MESSAGE, gateway.innerMessageHandler)
				client.HandleFunc(SYS_ON_CLOSE, gateway.innerCloseHandler)
				workers[value] = client.Dial()

				//加入一致性哈希
				globalConsistent.Add(value)
				gateway.reHashConnection()

				MyLog().Info("worker connected ", name, value)
			}
		}
		//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
	}
	return
}
