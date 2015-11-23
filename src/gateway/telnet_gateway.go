package main

import (
	"pillx"
	"fmt"
	"sync/atomic"
)

var workers map[uint32] *pillx.Response
var worker_id uint32 = 0;

var clients map[uint64] *pillx.Response
var chat_channel *pillx.Channel;

func innerConnectHandler(worker *pillx.Response, protocol pillx.IProtocol) {
	worker_id = atomic.AddUint32(&worker_id, 1)
	workers[worker_id] = worker
	fmt.Printf("worker %d 连接到此网关\n", worker_id)
}

func innerMessageHandler(worker *pillx.Response, protocol pillx.IProtocol) {
	//将gateway协议转化为客户端协议
	textProtocol := &pillx.TextProtocol{}
	fmt.Printf("test5\n")
	req := protocol.(*pillx.GateWayProtocol)
	textProtocol.Content = req.Content
	
	//发送给client
	//client := clients[req.Header.ClientId]
	fmt.Printf("发送给clients")
	//client.Send(textProtocol)
	chat_channel.Publish(textProtocol)
}

func innerCloseHandler(client *pillx.Response, protocol pillx.IProtocol) {

}

func outerConnectHandler(client *pillx.Response, protocol pillx.IProtocol) {
	clients[client.GetConn().Id] = client
	fmt.Printf("client %d 连接到此网关\n", client.GetConn().Id)
	
	//频道
	chat_channel.Subscribe(client);
}

func outerMessageHandler(client *pillx.Response, protocol pillx.IProtocol) {
	//将客户端协议转化为gateway协议
	gatewayProtocol := &pillx.GateWayProtocol{}
	header := &pillx.GatewayHeader{}
	gatewayProtocol.Header = header
	
	header.ClientId = client.GetConn().Id
	req := protocol.(*pillx.TextProtocol)
	header.Cmd = 0x0DDC
	header.Error = 0x0000
	header.Mark = 0xA8
	header.Size = uint16(len(req.Content))
	gatewayProtocol.Content = req.Content
	
	//发送给一个合适的worker
	worker := workers[worker_id]
	fmt.Printf("%x", gatewayProtocol.Header)
	fmt.Printf("%s", req.Content)
	worker.Send(gatewayProtocol)
	fmt.Printf("发送给worker %d\n", worker_id)
}

func outerCloseHandler(client *pillx.Response, protocol pillx.IProtocol) {
	delete(clients, client)
}

func main() {
	workers = make(map[uint32] *pillx.Response)
	clients = make(map[uint64] *pillx.Response)
	
	chat_channel = pillx.NewChannel("chat")
	
	innerServer := &pillx.Server{
		Addr:          ":10086",
		Handler:        nil,
		Protocol:		&pillx.GateWayProtocol{},
	}
	innerServer.HandleFunc(pillx.SYS_ON_CONNECT, innerConnectHandler)
	innerServer.HandleFunc(pillx.SYS_ON_MESSAGE, innerMessageHandler)
	innerServer.HandleFunc(pillx.SYS_ON_CLOSE, innerCloseHandler)
	fmt.Println("内部通信服务启动")
	go innerServer.ListenAndServe()
	
	outerServer := &pillx.Server{
		Addr:          ":8080",
		Handler:        nil,
		Protocol:		&pillx.TextProtocol{},
	}
	outerServer.HandleFunc(pillx.SYS_ON_CONNECT, outerConnectHandler)
	outerServer.HandleFunc(pillx.SYS_ON_MESSAGE, outerMessageHandler)
	outerServer.HandleFunc(pillx.SYS_ON_CLOSE, outerCloseHandler)
	fmt.Println("外部通信网关服务启动")
	outerServer.ListenAndServe()
}