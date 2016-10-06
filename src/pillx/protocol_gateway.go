package pillx

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type GatewayHeader struct {
	Mark     uint8
	Version  uint8
	Cmd      uint16
	Sid      uint32
	ClientId uint64
	Error    uint16
	Size     uint16
}

//将来以此为核心，讲gatewayprams作为一个数组，将获取到的值放入一个map中
type GatewayParamsHeader struct {
	Type uint8 //1.开始2.中断3.结束3.参数8.data
	Size uint16
	//NameSize uint8
}

type GateWayProtocol struct {
	Header *GatewayHeader
	//IPHeader *GatewayParamsHeader
	//IP       []byte
	//DataHeader *GatewayParamsHeader
	Content []byte
	//ParamsHeader     map[string]*GatewayParamsHeader
	//ParamsContent	map[string][]byte
}

const GATEWAY_VERSION uint8 = 0x01

func NewGatewayProtocol() (protocol *GateWayProtocol) {
	header := &GatewayHeader{
		Mark:     PROTO_HEADER_FIRSTCHAR,
		Version:  GATEWAY_VERSION,
		Sid:      0,
		ClientId: 0,
		Cmd:      0x0001,
		Error:    0,
		Size:     0,
	}
	gateway := &GateWayProtocol{
		Header: header,
	}
	return gateway
}

func (req *GateWayProtocol) New() (protocol IProtocol) {
	return new(GateWayProtocol)
}

func (gateway *GateWayProtocol) errorMsg(err_type uint8, err_num uint16, err_prama error) (err error) {
	//返回error
	errorMsg := &GatewayHeader{
		Mark:     PROTO_HEADER_FIRSTCHAR,
		Version:  GATEWAY_VERSION,
		Sid:      0,
		ClientId: 0,
		Cmd:      0x0001,
		Error:    err_num,
		Size:     0,
	}
	MyLog().Error(err_prama)
	errorBuf := new(bytes.Buffer)
	errBin := binary.Write(errorBuf, binary.BigEndian, errorMsg)
	if errBin != nil {
		fmt.Println("binary.Write failed:", errBin)
		return errBin
	}
	binary.Write(errorBuf, binary.BigEndian, 0)

	return &ProtocalError{
		err_type: err_type,
		err_msg:  errorBuf.Bytes(),
		err:      err_prama,
	}
}

func (gateway *GateWayProtocol) Analyze(client *Response) (err error) {
	header := new(GatewayHeader)
	gateway.Header = header
	buf := client.conn.buf

	if client.conn.connected_flg != true {
		client.conn.connected_flg = true
		//派发连接通知
		client.callbackServe(SYS_ON_CONNECT)
		return nil
	}
	//初始字节判断
	var mark_err error
	header.Mark, mark_err = buf.ReadByte()
	if mark_err != nil {
		return &ProtocalError{
			err_type: Protocal_Error_TYPE_DISCONNECT,
			err:      mark_err,
		}
	}
	//MyLog().Info(header.Mark)
	if header.Mark != PROTO_HEADER_FIRSTCHAR {
		//返回error
		return gateway.errorMsg(Protocal_Error_TYPE_COMMON, 0x0001, errors.New("request mark error"))
	}

	//取出版本号
	version, _ := buf.ReadByte()
	header.Version = version

	//取出cmd,size,error,全都是uint16,两个字节
	cmdB1, _ := buf.ReadByte()
	cmdB2, _ := buf.ReadByte()
	header.Cmd = uint16(cmdB1)<<8 | uint16(cmdB2)

	//取出请求id
	sId := make([]byte, 4)
	buf.Read(sId)
	header.Sid = binary.BigEndian.Uint32(sId)

	//连接id
	clientId := make([]byte, 8)
	buf.Read(clientId)
	b_buf := bytes.NewBuffer(clientId)
	binary.Read(b_buf, binary.BigEndian, &header.ClientId)

	errorB1, _ := buf.ReadByte()
	errorB2, _ := buf.ReadByte()
	header.Error = uint16(errorB1)<<8 | uint16(errorB2)

	//========相同协议处理======================================
	//paramsHeader := &GatewayParamsHeader{}
	//paramsHeader.Type, _ = buf.ReadByte()

	sizeB1, _ := buf.ReadByte()
	sizeB2, _ := buf.ReadByte()
	header.Size = uint16(sizeB1)<<8 | uint16(sizeB2)
	//gateway.DataHeader = paramsHeader

	/**
	//参数处理
	//取出ipsize
	ipParamsHeader := &GatewayParamsHeader{}
	ipsize, _ := buf.ReadByte()
	ip := make([]byte, int(ipsize))
	buf.Read(ip)
	ipParamsHeader.Size = ipsize
	gateway.IPHeader = ipParamsHeader
	gateway.IP = ip
	*/

	//fmt.Printf("%x", header)

	//根据size取出数据
	readNum := 0
	gateway.Content = make([]byte, header.Size)
	for readNum < int(header.Size) {
		readOnceNum, contentError := buf.Read(gateway.Content[readNum:])
		//fmt.Println(readOnceNum)
		if contentError != nil {
			return &ProtocalError{
				err_type: Protocal_Error_TYPE_DISCONNECT,
				err:      contentError,
			}
		}
		readNum += readOnceNum
	}
	//fmt.Printf("%s", gateway.Content)
	client.callbackServe(SYS_ON_MESSAGE)

	return nil
}

func (gateway *GateWayProtocol) Encode(msg interface{}) (buf []byte, err error) {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, msg.(*GateWayProtocol).Header)
	binary.Write(buff, binary.BigEndian, msg.(*GateWayProtocol).Content)
	//MyLog().Info(buff.Bytes())
	return buff.Bytes(), nil
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
