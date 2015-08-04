// new_server
package main

import (
	"io"
	"pillx"
	"fmt"
)

func helloHandler(client *pillx.Response, req *pillx.Request) {
	//fmt.Print(string(req.Content))
    io.WriteString(client, "World")
	
	channel := pillx.NewChannel("all")
	channel.Subscribe(client)
	channel.Publish(req)
}

func main() {
	server := &pillx.Server{
		Addr:          ":8080",
		Handler:        nil,
	}
	pillx.HandleFunc(0x0DDC, helloHandler)
	fmt.Println("pillX服务端引擎启动")
	server.ListenAndServe()
}
