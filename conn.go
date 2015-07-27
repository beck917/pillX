package pillx

import (
	"bufio"
	"io"
	"net"
	"sync"
	"errors"
	"fmt"
	"encoding/binary"
	"bytes"
)

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
	} else {
		response.conn.buf.Flush()
	}
	return
}

type RequestHeader struct {
	mark		uint8
	cmd 		uint16
	error		uint16
	size		uint16
}

type Request struct {
	Header		*RequestHeader
	Content		[]byte
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

func (c *Conn) errorResponse(response *Response, errorNum uint16) (err error) {
	//返回error
	errorMsg := &RequestHeader{
		mark:	0xA8,
		size:	0,
		cmd:	0x0001,
		error:	errorNum,
	}
	errorBuf := new(bytes.Buffer)
	errBin := binary.Write(errorBuf, binary.BigEndian, errorMsg)
	if (errBin != nil) {
		fmt.Println("binary.Write failed:", errBin)
		return errBin
	}
	fmt.Printf("%x\n", errorBuf.Bytes())
	response.Write(errorBuf.Bytes())
	return nil
}

func (c *Conn) readRequest() (response *Response, err error) {
	var req *Request
	req = new(Request)
	reqHeader := new(RequestHeader)
	
	response = &Response{
		conn:          c,
		req:           req,
		contentLength: -1,
	}
	
	//初始字节判断
	var mark_err error
	reqHeader.mark, mark_err = c.buf.ReadByte()
	if (mark_err != nil) {
		response.conn.remonte_conn.Close()
		return nil, mark_err
	}
	if (reqHeader.mark != 0xa8) {
		c.buf.Reader.Reset(c.lr)
		//返回error
		c.errorResponse(response, 0x0001)
		return nil, errors.New("request mark error")
	}
	//取出cmd,size,error,全都是uint16,两个字节
	cmdB1, _ := c.buf.ReadByte()
	cmdB2, _ := c.buf.ReadByte()
	reqHeader.cmd = uint16(cmdB1) << 8 | uint16(cmdB2)
	
	errorB1, _ := c.buf.ReadByte()
	errorB2, _ := c.buf.ReadByte()
	reqHeader.error = uint16(errorB1) << 8 | uint16(errorB2)
	
	sizeB1, _ := c.buf.ReadByte()
	sizeB2, _ := c.buf.ReadByte()
	reqHeader.size = uint16(sizeB1) << 8 | uint16(sizeB2)
	
	/*
	reqBuf := c.buf.Read(make([]byte, 7))
	Request(reqBuf).size
	*/
	//根据size取出数据
	readNum := 0
	req.Content = make([]byte, reqHeader.size)
	for readNum < int(reqHeader.size) {
		readOnceNum,contentError := c.buf.Read(req.Content[readNum:])
		if contentError != nil {
			c.buf.Reader.Reset(c.lr)
			c.errorResponse(response, 0x0002)
			return nil, errors.New("request size error")
		}
		readNum += readOnceNum
	}

	req.Header = reqHeader
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
	}
}