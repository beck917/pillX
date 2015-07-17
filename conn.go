package pillx

import (
	"bufio"
	"io"
	"net"
	"sync"
)

type Request struct {
	mark		uint8
	size		uint16
	cmd 		uint16
	content		[]byte
}

// A response represents the server side of aresponse.
type Response struct {
	conn          *Conn
	req           *Request // request for this response

	written       int64 // number of bytes written in body
	contentLength int64 // explicitly-declared Content-Length; or -1
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
	}
	return
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
	var req *Request
	req = new(Request)//c.buf.Reader
	req.mark = c.buf.ReadByte()
	if (req.mark != 0xa8) {
		
	}
	//req.name = c.buf.Read()
	req.name = 0x0DDC
	req.conn = c
	req.Buf = c.buf
	
	response = &Response{
		conn:          c,
		req:           req,
		contentLength: -1,
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
		
		ServerHandler{c.server}.serve(w, w.req)
		break
	}
}