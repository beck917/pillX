package pillx

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type PillProtocolHeader struct {
	Mark    uint8
	Version uint8
	Cmd     uint16
	Sid     uint32
	Error   uint16
	Size    uint16
}

const PILL_VERSION uint8 = 0x01

type PillProtocol struct {
	Header  *PillProtocolHeader
	Content []byte
}

func NewPillProtocol() (protocol *PillProtocol) {
	header := &PillProtocolHeader{
		Mark:    PROTO_HEADER_FIRSTCHAR,
		Version: GATEWAY_VERSION,
		Sid:     0,
		Cmd:     0x0001,
		Error:   0,
		Size:    0,
	}
	gateway := &PillProtocol{
		Header: header,
	}
	return gateway
}

func (req *PillProtocol) New() (protocol IProtocol) {
	return new(PillProtocol)
}

func (req *PillProtocol) errorMsg(err_type uint8, err_num uint16, err_prama error) (err error) {
	//返回error
	errorMsg := &PillProtocolHeader{
		Mark:    PROTO_HEADER_FIRSTCHAR,
		Size:    0,
		Version: PILL_VERSION,
		Sid:     0,
		Cmd:     0x0001,
		Error:   err_num,
	}
	errorBuf := new(bytes.Buffer)
	errBin := binary.Write(errorBuf, binary.BigEndian, errorMsg)
	if errBin != nil {
		fmt.Println("binary.Write failed:", errBin)
		return errBin
	}

	return &ProtocalError{
		err_type: err_type,
		err_msg:  errorBuf.Bytes(),
		err:      err_prama,
	}
}

func (req *PillProtocol) Analyze(client *Response) (err error) {
	buf := client.conn.buf
	reqHeader := new(PillProtocolHeader)
	req.Header = reqHeader

	if client.conn.connected_flg != true {
		client.conn.connected_flg = true
		//派发连接通知
		client.callbackServe(SYS_ON_CONNECT)
		return nil
	}

	//初始字节判断
	var mark_err error
	reqHeader.Mark, mark_err = buf.ReadByte()
	if mark_err != nil {
		return &ProtocalError{
			err_type: Protocal_Error_TYPE_DISCONNECT,
			err:      mark_err,
		}
	}

	if reqHeader.Mark != PROTO_HEADER_FIRSTCHAR {
		//返回error
		return req.errorMsg(Protocal_Error_TYPE_COMMON, 0x0001, errors.New("request mark error"))
	}

	//取出版本号
	version, _ := buf.ReadByte()
	reqHeader.Version = version

	//取出cmd,size,error,全都是uint16,两个字节
	cmdB1, _ := buf.ReadByte()
	cmdB2, _ := buf.ReadByte()
	reqHeader.Cmd = uint16(cmdB1)<<8 | uint16(cmdB2)

	//取出请求id
	sId := make([]byte, 4)
	buf.Read(sId)
	reqHeader.Sid = binary.BigEndian.Uint32(sId)

	errorB1, _ := buf.ReadByte()
	errorB2, _ := buf.ReadByte()
	reqHeader.Error = uint16(errorB1)<<8 | uint16(errorB2)

	sizeB1, _ := buf.ReadByte()
	sizeB2, _ := buf.ReadByte()
	reqHeader.Size = uint16(sizeB1)<<8 | uint16(sizeB2)
	//fmt.Printf("%x", reqHeader)
	//根据size取出数据
	readNum := 0
	req.Content = make([]byte, reqHeader.Size)
	for readNum < int(reqHeader.Size) {
		readOnceNum, contentError := buf.Read(req.Content[readNum:])
		//fmt.Println(readOnceNum)
		if contentError != nil {
			return &ProtocalError{
				err_type: Protocal_Error_TYPE_DISCONNECT,
				err:      contentError,
			}
		}
		readNum += readOnceNum
	}
	//fmt.Printf("%s", req.Content)
	client.callbackServe(SYS_ON_MESSAGE)

	return nil
}

func (req *PillProtocol) Encode(msg interface{}) (buf []byte, err error) {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, msg.(*PillProtocol).Header)
	binary.Write(buff, binary.BigEndian, msg.(*PillProtocol).Content)

	return buff.Bytes(), nil
}

func (req *PillProtocol) Decode(buf []byte) (err error) {
	return nil
}

func (req *PillProtocol) SetCmd(cmd uint16) {
	req.Header.Cmd = cmd
}

func (req *PillProtocol) GetCmd() (cmd uint16) {
	return req.Header.Cmd
}
