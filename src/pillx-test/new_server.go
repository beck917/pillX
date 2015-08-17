// new_server
package main

import (
	//"io"
	"pillx"
	"fmt"
)

func helloHandler(client *pillx.Response, req pillx.IProtocol) {
	//fmt.Print(string(req.Content))
    //
	client.Send(req)
	//io.WriteString(client, "World")
	
	channel := pillx.NewChannel("all")
	channel.Subscribe(client)
	channel.Publish(req)
}

func closeHandler(client *pillx.Response, req pillx.IProtocol) {
	fmt.Print("closed")
}

func main() {
	server := &pillx.Server{
		Addr:          ":8080",
		Handler:        nil,
		Protocol:		new(pillx.PillProtocol),
	}
	server.HandleFunc(0x0DDC, helloHandler)
	server.HandleFunc(pillx.SYS_ON_CLOSE, closeHandler)
	fmt.Println("pillX服务端引擎启动")
	server.ListenAndServe()
}
