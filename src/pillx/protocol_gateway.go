package pillx

import (
	"errors"
	"bytes"
	"fmt"
	"encoding/binary"
)

type GatewayHeader struct {
	Mark		uint8
	Cmd 		uint16
	ClientId	uint64
	Error		uint16
	Size		uint16
}

type GateWayProtocol struct {
	Header			*GatewayHeader
	Content			[]byte
}

func (gateway *GateWayProtocol) New() (protocol IProtocol) {
	return new(GateWayProtocol)
}

func (gateway *GateWayProtocol) errorMsg(err_type uint8, err_num uint16, err_prama error) (err error) {
	//返回error
	errorMsg := &PillProtocolHeader{
		Mark:	0xA8,
		Size:	0,
		Cmd:	0x0001,
		Error:	err_num,
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

func (gateway *GateWayProtocol) Analyze(client *Response) (err error) {
	header := new(GatewayHeader)
	gateway.Header = header
	buf := client.conn.buf

	if (client.conn.handshake_flg != true) {
		client.conn.handshake_flg = true
		//派发连接通知
		client.callbackServe(SYS_ON_CONNECT)
		return nil
	}
	//初始字节判断
	var mark_err error
	header.Mark, mark_err = buf.ReadByte()
	if (mark_err != nil) {
		client.callbackServe(SYS_ON_CLOSE)
		return &ProtocalError{
					err_type:	Protocal_Error_TYPE_DISCONNECT,
					err:		mark_err,
				}
	}
	
	if (header.Mark != 0xa8) {
		//返回error
		return gateway.errorMsg(Protocal_Error_TYPE_COMMON, 0x0001, errors.New("request mark error"))
	}
	fmt.Println("test1")
	//取出cmd,size,error,全都是uint16,两个字节
	cmdB1, _ := buf.ReadByte()
	cmdB2, _ := buf.ReadByte()
	header.Cmd = uint16(cmdB1) << 8 | uint16(cmdB2)
	
	clientId := make([]byte, 8)
	buf.Read(clientId)
	b_buf := bytes.NewBuffer(clientId)
    binary.Read(b_buf, binary.BigEndian, &header.ClientId)

	errorB1, _ := buf.ReadByte()
	errorB2, _ := buf.ReadByte()
	header.Error = uint16(errorB1) << 8 | uint16(errorB2)
	
	sizeB1, _ := buf.ReadByte()
	sizeB2, _ := buf.ReadByte()
	header.Size = uint16(sizeB1) << 8 | uint16(sizeB2)
	
	//根据size取出数据
	readNum := 0
	gateway.Content = make([]byte, header.Size)
	for readNum < int(header.Size) {
		readOnceNum,contentError := buf.Read(gateway.Content[readNum:])
		if contentError != nil {
			return &ProtocalError{
						err_type:	Protocal_Error_TYPE_DISCONNECT,
						err:		contentError,
					}
		}
		readNum += readOnceNum
	}
	client.callbackServe(SYS_ON_MESSAGE)
	
	return nil
}

func (gateway *GateWayProtocol) Encode(msg interface{}) (buf []byte, err error) {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, msg.(*GateWayProtocol).Header)
	binary.Write(buff, binary.BigEndian, msg.(*GateWayProtocol).Content)

	return buff.Bytes(),nil
}

func (gateway *GateWayProtocol) Decode(buf []byte) (err error) {
	return nil
}

func (gateway *GateWayProtocol) SetCmd(cmd uint16) {
	gateway.Header.Cmd = cmd
}

func (gateway *GateWayProtocol) GetCmd() (cmd uint16) {
	return gateway.Header.Cmd
}