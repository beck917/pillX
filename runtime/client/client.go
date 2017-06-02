package main

import (
	"fmt"
	"pillx"
	//"strconv"
	"time"
)

func helloHandler(client *pillx.Response, protocol pillx.IProtocol) {
	fmt.Println("test3")
	req := protocol.(*pillx.PillProtocol)
	fmt.Printf("%x", req.Header)
	fmt.Printf("%s", req.Content)
	//req.Content = []byte("user" + strconv.FormatUint(uint64(req.Header.Sid), 10) + ":  " + string(req.Content) + "\r\n")
	//req.Header.Size = uint16(len(req.Content))
	//client.Send(req)
}

func main() {
	//连接gateway服务器
	client := &pillx.Server{
		Addr:     "115.29.164.163:9127",
		Handler:  nil,
		Protocol: &pillx.PillProtocol{},
	}
	fmt.Println("内部通信服务启动")
	//client.HandleFunc(pillx.SYS_ON_CONNECT, connectHandler)
	client.HandleFunc(0x0DDC, helloHandler)
	//client.HandleFunc(pillx.SYS_ON_MESSAGE, messageHandler)
	rc := client.Dial()
	fmt.Println("连接")

	//发送数据
	pill := &pillx.PillProtocol{}
	pillHeader := &pillx.PillProtocolHeader{
		Mark:    pillx.PROTO_HEADER_FIRSTCHAR,
		Size:    4,
		Version: pillx.PILL_VERSION,
		Sid:     110,
		Cmd:     0x0DDC,
		Error:   0,
	}
	pill.Header = pillHeader
	pill.Content = []byte("test")

	for i := range []int{1, 2, 3} {
		fmt.Println(i)
		rc.Send(pill)
		time.Sleep(time.Second * 5)
	}
	rc.Send(pill)
	for {
		time.Sleep(1 * 1e9)
	}
}
