package main

import (
	"pillx"
	"fmt"
)

var workers map[int] *pillx.Response

func innerConnectHandler(client *pillx.Response, req pillx.IProtocol) {

}

func innerMessageHandler(client *pillx.Response, req pillx.IProtocol) {

}

func innerCloseHandler(client *pillx.Response, req pillx.IProtocol) {

}

func outerConnectHandler(client *pillx.Response, req pillx.IProtocol) {

}

func outerMessageHandler(client *pillx.Response, req pillx.IProtocol) {

}

func outerCloseHandler(client *pillx.Response, req pillx.IProtocol) {

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
		Protocol:		&pillx.GateWayProtocol{},
	}
	outerServer.HandleFunc(pillx.SYS_ON_CONNECT, outerConnectHandler)
	outerServer.HandleFunc(pillx.SYS_ON_MESSAGE, outerMessageHandler)
	outerServer.HandleFunc(pillx.SYS_ON_CLOSE, outerCloseHandler)
	fmt.Println("外部通信网关服务启动")
	outerServer.ListenAndServe()
}