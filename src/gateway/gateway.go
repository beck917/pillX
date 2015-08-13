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
	client := clients[worker.Id]
	client.Send(pillProtocol)
}

func innerCloseHandler(client *pillx.Response, protocol pillx.IProtocol) {

}

func outerConnectHandler(client *pillx.Response, protocol pillx.IProtocol) {
	clients[client.Id] = client
}

func outerMessageHandler(client *pillx.Response, protocol pillx.IProtocol) {
	//将客户端协议转化为gateway协议
	gatewayProtocol := &pillx.GateWayProtocol{}
	header := &pillx.GatewayHeader{}
	gatewayProtocol.Header = header
	
	header.ClientId = client.Id
	req := protocol.(*pillx.PillProtocol)
	header.Cmd = req.Header.Cmd
	header.Error = req.Header.Error
	header.Mark = req.Header.Mark
	header.Size = req.Header.Size
	
	buf, _ := gatewayProtocol.Encode(gatewayProtocol)
	gatewayProtocol.Content = buf
	
	//发送给一个合适的worker
	worker := workers[worker_id]
	worker.Send(gatewayProtocol)
}

func outerCloseHandler(client *pillx.Response, protocol pillx.IProtocol) {

}

func main() {
	innerServer := &pillx.Server{
		Addr:          ":10086",
		Handler:        nil,
		Protocol:		&pillx.GateWayProtocol{},
	}
	innerServer.HandleFunc(pillx.SYS_ON_CONNECT, innerConnectHandler)
	innerServer.HandleFunc(pillx.SYS_ON_MESSAGE, innerMessageHandler)
	innerServer.HandleFunc(pillx.SYS_ON_CLOSE, innerCloseHandler)
	fmt.Println("内部通信服务启动")
	innerServer.ListenAndServe()
	
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