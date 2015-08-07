package pillx

import (
	//"fmt"
	"bufio"
	"io"
	"net"
	"sync"
)

type IProtocol interface {
	Analyze(buf *bufio.ReadWriter) (err error)
	Encode(msg interface{}) (buf []byte, err error)
	Decode(buf []byte) (err error)
	GetCmd() (cmd uint16)
}

const Protocal_Error_TYPE_COMMON uint8 = 1
const Protocal_Error_TYPE_DISCONNECT uint8 = 2

type ProtocalError struct {
	err_type uint8
	err_msg []byte
	err error
}

func (e *ProtocalError) Error() string{
	return e.err.Error()
}

// A response represents the server side of aresponse.
type Response struct {
	conn          *Conn
	protocol      IProtocol // request for this response
	channels	  map[string]*Channel
}

//取消所有订阅频道
func (response *Response) unSubscribeAllChannels() (err error) {
	for _,channel :=range response.channels {
		channel.UnSubscribe(response)
	}
	return nil
}

func (response *Response) Write(data []byte) (n int, err error) {
	return response.write(len(data), data, "")
}

func (response *Response) write(lenData int, dataB []byte, dataS string) (n int, err error) {
	if dataB != nil {
		n, err = response.conn.buf.Write(dataB)
	} else {
		n, err = response.conn.buf.WriteString(dataS)
	}

	if err != nil {
		response.conn.remonte_conn.Close()
	} else {
		response.conn.buf.Flush()
	}
	return
}

func (response *Response) SendContent(content []byte) {
	
}

func (response *Response) Send(msg interface{}) {
	buf, _ := response.protocol.Encode(msg)
	
	response.Write(buf)
}

// A conn represents the server side of connection.
type Conn struct {
	remote_addr 		string
	server 				*Server
	remonte_conn 		net.Conn
	io_writer			io.Writer
	io_writer_err       error                // any errors writing to w
	sr         			liveSwitchReader     // where the LimitReader reads from; usually the rwc
	lr         			*io.LimitedReader    // io.LimitReader(sr)
	buf 				*bufio.ReadWriter
	
	mu 			sync.Mutex
}

func (c *Conn) readRequest() (response *Response, err error) {
	protocol := c.server.Protocol
	
	response = &Response{
		conn:          c,
		protocol:      protocol,
		channels: 	   make(map[string]*Channel),
	}
	
	err = protocol.Analyze(c.buf)
	if err != nil {
		switch err.(type) {
			case *ProtocalError: 
				switch err.(*ProtocalError).err_type {
					case Protocal_Error_TYPE_DISCONNECT:
						//取消订阅频道
						response.unSubscribeAllChannels()
						
						c.buf.Reader.Reset(c.lr)
						response.conn.remonte_conn.Close()
						
						//发送onclose通知
						ServerHandler{c.server}.serve(response, protocol)
						return nil, err
						break
					case Protocal_Error_TYPE_COMMON:
						//数据重置
						c.buf.Reader.Reset(c.lr)
						response.Write(err.(*ProtocalError).err_msg)
						break
					default:
						return nil, err
				}
			default:
				return nil, nil
		}
		
	}
	return response, nil
}

// Serve a new connection.
func (c *Conn) serve() {
	//origConn := c.rwc // copy it before it's set nil on Close or Hijack
	for {
		w, err := c.readRequest()
		
		if err != nil {
			break
		}
		
		ServerHandler{c.server}.serve(w, w.protocol)
	}
}