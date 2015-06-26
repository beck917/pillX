// new_server
package main

import (
	"fmt"
	"../../pillX"
)

func helloHandler(rw pillx.ResponseWriter, req *pillx.Request) {
    io.WriteString(rw, "hello world")
}


func main() {
	server := &pillx.Server{
		Addr:           ":8080",
		Handler:        nil,
	}
	pillx.HandleFunc(0x0DDC, helloHandler)
	server.ListenAndServe()
}
