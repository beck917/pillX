package pillx

import (
	//"errors"
	"bytes"
	//"fmt"
	"encoding/binary"
)

type GatewayHeader struct {
	Mark		uint8
	Cmd 		uint16
	Error		uint16
	Size		uint16
	ClientId	uint64
}

type GateWayProtocol struct {
	Header			*GatewayHeader
	Content			[]byte
}

func (gateway *GateWayProtocol) New() (protocol IProtocol) {
	return new(GateWayProtocol)
}

func (gateway *GateWayProtocol) Analyze(client *Response) (err error) {
	if (client.handshake_flg != true) {
		client.handshake_flg = true
		//派发连接通知
		client.callbackServe(SYS_ON_CONNECT)
		return nil
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