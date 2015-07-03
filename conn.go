package pillx

import (
	"bufio"
	"io"
	"net"
	"sync"
)

type Request struct {
	name uint16
}

// A response represents the server side of an HTTP response.
type Response struct {
	conn          *Conn
	req           *Request // request for this response

	written       int64 // number of bytes written in body
	contentLength int64 // explicitly-declared Content-Length; or -1

	// close connection after this reply.  set on request and
	// updated after response from handler if there's a
	// "Connection: keep-alive" response header and a
	// Content-Length.
	closeAfterReply bool
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
		//cw.res.conn.rwc.Close()
	}
	return
}

// A conn represents the server side of connection.
type Conn struct {
	remote_addr 		string
	server 				*Server
	remonte_conn 		net.Conn
	io_writer			io.Writer
	io_writer_err       			error                // any errors writing to w
	sr         			liveSwitchReader     // where the LimitReader reads from; usually the rwc
	lr         			*io.LimitedReader    // io.LimitReader(sr)
	buf 				*bufio.ReadWriter
	
	mu 			sync.Mutex
}

func (c *Conn) readRequest() (response *Response, err error) {
	var req *Request
	req = new(Request)//c.buf.Reader
	//req.name = c.buf.Read();
	req.name = 0x0DDC;
	
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