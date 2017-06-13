package controllers

import (
	"time"
	//"fmt"
	"application/libraries/helpers"

	"github.com/beck917/pillX/Proto"
	"github.com/beck917/pillX/pillx"

	log "github.com/Sirupsen/logrus"
	"github.com/bitly/go-simplejson"
	"github.com/golang/protobuf/proto"
)

func OnMessageHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)
	jsonObj, jsonErr := simplejson.NewJson(req.Content)
	if jsonErr != nil {
		//记录错误
		return
	}

	method, err := jsonObj.Get("method").String()

	if err != nil {
		return
	}
	cmd := helpers.Crc16([]byte(method))
	req.SetCmd(cmd)

	helpers.GlobalWorker.(*pillx.Worker).InnerServer.Handler.Serve(client, protocol)
}

func IndexHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)

	//解析content
	//数据处理到proto类
	chatData := &Proto.ChatData{}
	err := proto.Unmarshal(req.Content, chatData)

	if err != nil {
		pillx.MyLog().WithField("data", string(req.Content)).Error(err)
	} else {
		var ip string
		if _, ok := pillx.WorkerClients[req.Header.ClientId]; ok {
			ip = pillx.WorkerClients[req.Header.ClientId].IP
		} else {
			ip = "0.0.0.0"
		}
		pillx.MyLog().WithFields(log.Fields{
			"ip":      ip,
			"byte":    string(req.Content),
			"uid":     chatData.GetUid(),
			"msg":     chatData.GetMsg(),
			"msgjson": chatData.GetMsgjson(),
			//"roomid":  chatData.Header.GetRoomId(),
			"type": "chat",
		}).Info("处理消息内容")
	}
	chatData.Timestamp = proto.Int32(int32(time.Now().Unix()))
	req.Content, _ = proto.Marshal(chatData)
	req.Header.Size = uint16(len(req.Content))

	//发送到所有网关
	pillx.SendAllGateWay(pillx.GatewayPools, req)
}

func LoginHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)

	//解析content
	_, jsonErr := simplejson.NewJson(req.Content)
	if jsonErr != nil {
		//记录错误
		return
	}

	jsonRet, _ := simplejson.NewJson([]byte(`{}`))
	jsonRet.Set("method", "loginresult")
	jsonRet.Set("msg", "loginok")
	jsonRet.Set("result", "valid")
	jsonRet.Set("userid", 10001)
	req.Content, _ = jsonRet.Encode()
	req.Header.Size = uint16(len(req.Content))

	//发送到所有网关
	pillx.SendAllGateWay(pillx.GatewayPools, req)
}
