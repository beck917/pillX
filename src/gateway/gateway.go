package main

import (
	"pillx"
	"fmt"
	"sync/atomic"
)

var workers map[uint32] *pillx.Response
var worker_id uint32 = 0;

var clients map[uint64] *pillx.Response

func innerConnectHandler(worker *pillx.Response, protocol pillx.IProtocol) {
	worker_id = atomic.AddUint32(&worker_id, 1)
	workers[worker_id] = worker
	fmt.Printf("worker %d 连接到此网关\n", worker_id)
}

func innerMessageHandler(worker *pillx.Response, protocol pillx.IProtocol) {
	//将gateway协议转化为客户端协议
	pillProtocol := &pillx.PillProtocol{}
	header := &pillx.PillProtocolHeader{}
	pillProtocol.Header = header
	
	req := protocol.(*pillx.GateWayProtocol)
	header.Cmd = req.Header.Cmd
	header.Error = req.Header.Error
	header.Mark = req.Header.Mark
	header.Size = req.Header.Size
	
	buf, _ := pillProtocol.Encode(pillProtocol)
	pillProtocol.Content = buf
	
	//发送给client
	client := clients[worker.GetConn().Id]
	client.Send(pillProtocol)
	fmt.Printf("发送给client %d\n", worker.GetConn().Id)
}

func innerCloseHandler(client *pillx.Response, protocol pillx.IProtocol) {

}

func outerConnectHandler(client *pillx.Response, protocol pillx.IProtocol) {
	clients[client.GetConn().Id] = client
	fmt.Printf("client %d 连接到此网关\n", client.GetConn().Id)
}

func outerMessageHandler(client *pillx.Response, protocol pillx.IProtocol) {
	//将客户端协议转化为gateway协议
	gatewayProtocol := &pillx.GateWayProtocol{}
	header := &pillx.GatewayHeader{}
	gatewayProtocol.Header = header
	
	header.ClientId = client.GetConn().Id
	req := protocol.(*pillx.PillProtocol)
	header.Cmd = req.Header.Cmd//cmd被覆盖的问题
	header.Error = req.Header.Error
	header.Mark = req.Header.Mark
	header.Size = req.Header.Size
	gatewayProtocol.Content = req.Content

	//buf, _ := gatewayProtocol.Encode(gatewayProtocol)
	
	//发送给一个合适的worker
	worker := workers[worker_id]
	worker.Send(gatewayProtocol)
	fmt.Printf("发送给worker %d\n", worker_id)
}

func outerCloseHandler(client *pillx.Response, protocol pillx.IProtocol) {

}

func main() {
	workers = make(map[uint32] *pillx.Response)
	clients = make(map[uint64] *pillx.Response)
	
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
		Protocol:		&pillx.PillProtocol{},
	}
	outerServer.HandleFunc(pillx.SYS_ON_CONNECT, outerConnectHandler)
	outerServer.HandleFunc(pillx.SYS_ON_MESSAGE, outerMessageHandler)
	outerServer.HandleFunc(pillx.SYS_ON_CLOSE, outerCloseHandler)
	fmt.Println("外部通信网关服务启动")
	outerServer.ListenAndServe()
}