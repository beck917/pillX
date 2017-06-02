package pillx

import (
	"fmt"
	"net"
	"os"

	log "github.com/Sirupsen/logrus"
)

type Client struct {
	Addr string
}

func (client *Client) Dail() (p Pool, err error) {
	// create a factory() to be used with channel based pool

	factory := func() (*PoolConn, error) {
		server := &Server{
			Protocol: &GateWayProtocol{},
		}
		rc, err := net.Dial("tcp", client.Addr)
		if err != nil {
			log.WithError(err).Error(os.Stderr, "Fatal error: ", err.Error())
		}
		c, _ := server.newConn(rc)
		response := &Response{
			conn:     c,
			protocol: server.Protocol,
		}

		poolconn := &PoolConn{
			response: response,
		}
		return poolconn, err
	}

	/**
	factory := func() (net.Conn, error) {
		rc, err := net.Dial("tcp", client.Addr)
		if err != nil {
			log.WithError(err).Error(os.Stderr, "Fatal error: %s", err.Error())
		}
		return rc, err
	}
	*/

	// create a new channel based pool with an initial capacity of 5 and maximum
	// capacity of 30. The factory will create 5 initial connections and put it
	// into the pool.
	p, err = NewChannelPool(1000, 5000, factory)

	return p, err

	// now you can get a connection from the pool, if there is no connection
	// available it will create a new one via the factory function.
	//conn, err := p.Get()
}

func (client *Server) Dial() (response *Response) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", client.Addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	rc, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	//defer rc.Close()
	c, err := client.newConn(rc)
	response = &Response{
		conn:     c,
		protocol: client.Protocol,
	}
	go c.clientServe()
	return response
}
