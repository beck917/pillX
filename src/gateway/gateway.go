package main

import (
	"pillx"
)

var workers map[int] *pillx.Response

func connectHandler(client *pillx.Response, req *pillx.Request) {
	
}

func main() {
	server := &pillx.Server{
		Addr:          ":8080",
		Handler:        nil,
	}
	pillx.HandleFunc(0x0DDC, helloHandler)
	pillx.HandleFunc(0x0003, connectHandler)
	fmt.Println("pillX服务端引擎启动")
	server.ListenAndServe()
}