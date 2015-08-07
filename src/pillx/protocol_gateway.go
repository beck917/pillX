package pillx

import (
	"bufio"
	//"errors"
	"bytes"
	//"fmt"
)

type GateWayProtocol struct {
	Content			[]byte
	handshake_flg	bool
	cmd				uint16
}

func (gateway *GateWayProtocol) Analyze(buf *bufio.ReadWriter) (err error) {
	if (gateway.handshake_flg != true) {
		gateway.handshake_flg = true
		//派发
		gateway.cmd = SYS_ON_CONNECT
		return nil
	}
	
	gateway.cmd = SYS_ON_MESSAGE
	return nil
}

func (gateway *GateWayProtocol) Encode(msg interface{}) (buf []byte, err error) {
	buff := new(bytes.Buffer)

	return buff.Bytes(),nil
}

func (gateway *GateWayProtocol) Decode(buf []byte) (err error) {
	return nil
}

func (gateway *GateWayProtocol) GetCmd() (cmd uint16) {
	return gateway.cmd
}