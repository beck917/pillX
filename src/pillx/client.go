package pillx

import (
	"os"
	"fmt"
	"net"
)

type Client struct {
	Server
}

func (clinet *Server) Dial() {
    tcpAddr, err := net.ResolveTCPAddr("tcp4", clinet.Addr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
 
    rc, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
	defer rc.Close()
	c, err := clinet.newConn(rc)
	c.serve()
}