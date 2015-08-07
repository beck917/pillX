package pillx

import (
	"errors"
	"encoding/binary"
	"bytes"
	"fmt"
)

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

func (req *Request) New() (protocol IProtocol) {
	return new(Request)
}

func (req *Request) errorMsg(err_type uint8, err_num uint16, err_prama error) (err error) {
	//返回error
	errorMsg := &RequestHeader{
		mark:	0xA8,
		size:	0,
		cmd:	0x0001,
		error:	err_num,
	}
	errorBuf := new(bytes.Buffer)
	errBin := binary.Write(errorBuf, binary.BigEndian, errorMsg)
	if (errBin != nil) {
		fmt.Println("binary.Write failed:", errBin)
		return errBin
	}
	
	return &ProtocalError{
		err_type:	err_type,
		err_msg:	errorBuf.Bytes(),
		err:		err_prama,
	}
}

func (req *Request) Analyze(client *Response) (err error) {
	buf := client.conn.buf
	reqHeader := new(RequestHeader)
	req.Header = reqHeader
	
	//初始字节判断
	var mark_err error
	reqHeader.mark, mark_err = buf.ReadByte()
	if (mark_err != nil) {
		client.callbackServe(SYS_ON_CLOSE)
		return &ProtocalError{
					err_type:	Protocal_Error_TYPE_DISCONNECT,
					err:		mark_err,
				}
	}
	if (reqHeader.mark != 0xa8) {
		//返回error
		return req.errorMsg(Protocal_Error_TYPE_COMMON, 0x0001, errors.New("request mark error"))
	}
	
	//取出cmd,size,error,全都是uint16,两个字节
	cmdB1, _ := buf.ReadByte()
	cmdB2, _ := buf.ReadByte()
	reqHeader.cmd = uint16(cmdB1) << 8 | uint16(cmdB2)
	
	errorB1, _ := buf.ReadByte()
	errorB2, _ := buf.ReadByte()
	reqHeader.error = uint16(errorB1) << 8 | uint16(errorB2)
	
	sizeB1, _ := buf.ReadByte()
	sizeB2, _ := buf.ReadByte()
	reqHeader.size = uint16(sizeB1) << 8 | uint16(sizeB2)
	
	//根据size取出数据
	readNum := 0
	req.Content = make([]byte, reqHeader.size)
	for readNum < int(reqHeader.size) {
		readOnceNum,contentError := buf.Read(req.Content[readNum:])
		if contentError != nil {
			return &ProtocalError{
						err_type:	Protocal_Error_TYPE_DISCONNECT,
						err:		contentError,
					}	
			//return req.errorMsg(Protocal_Error_TYPE_COMMON, 0x0002, errors.New("request size error"))
		}
		readNum += readOnceNum
	}
	
	return nil
}

func (req *Request) Encode(msg interface{}) (buf []byte, err error) {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, msg.(*Request).Header)
	binary.Write(buff, binary.BigEndian, msg.(*Request).Content)

	return buff.Bytes(),nil
}

func (req *Request) Decode(buf []byte) (err error) {
	return nil
}

func (req *Request) SetCmd(cmd uint16) {
	req.Header.cmd = cmd
}

func (req *Request) GetCmd() (cmd uint16) {
	return req.Header.cmd
}