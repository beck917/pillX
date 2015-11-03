package main

import (
	"pillx"
	"fmt"
	"strconv"
)

func connectHandler(worker *pillx.Response, protocol pillx.IProtocol) {
	fmt.Println("连接gateway成功")
}

func helloHandler(client *pillx.Response, protocol pillx.IProtocol) {
	fmt.Println("test3")
	req := protocol.(*pillx.GateWayProtocol)
	fmt.Printf("%x", req.Header)
	fmt.Printf("%s", req.Content)
	req.Content =  []byte("user" + strconv.FormatUint(req.Header.ClientId, 10) + ":  " + string(req.Content)+"\r\n")
	req.Header.Size = uint16(len(req.Content))
	client.Send(req)
}

func messageHandler(worker *pillx.Response, protocol pillx.IProtocol) {
	fmt.Println("test4")
}

func main() {
	//连接gateway服务器
	client := &pillx.Server{
		Addr:          "127.0.0.1:10086",
		Handler:        nil,
		Protocol:		&pillx.GateWayProtocol{},
	}
	fmt.Println("内部通信服务启动")
	client.HandleFunc(pillx.SYS_ON_CONNECT, connectHandler)
	client.HandleFunc(0x0DDC, helloHandler)
	client.HandleFunc(pillx.SYS_ON_MESSAGE, messageHandler)
	client.Dial()
}