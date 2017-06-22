package pillx

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"strings"
)

type WebSocketHeader struct {
	OpcodeByte  byte
	PayloadByte byte
}

type WebSocketProtocol struct {
	Header  *WebSocketHeader
	Content []byte
}

func (this *WebSocketProtocol) New() (protocol IProtocol) {
	return new(WebSocketProtocol)
}

func (websocket *WebSocketProtocol) Analyze(client *Response) (err error) {
	//client.conn.mu.Lock()
	//defer client.conn.mu.Unlock()

	if client.conn.connected_flg != true {
		client.conn.connected_flg = true
		//派发连接通知
		client.callbackServe(SYS_ON_CONNECT)
		return nil
	}

	_, headerReq := websocket.Handshake(client)

	if len(headerReq) != 0 {
		//返回数据
		client.conn.buf.Write(headerReq)
		client.conn.buf.Flush()
		return nil
	}

	var (
		opcode byte
	)
	header := new(WebSocketHeader)
	websocket.Header = header
	buf := client.conn.buf

	//读取fin
	header.OpcodeByte, err = buf.ReadByte()
	if err != nil {
		//返回error
		return &ProtocalError{
			err_type: Protocal_Error_TYPE_DISCONNECT,
			err:      errors.New("EOF closed"),
		}
	}
	fin := header.OpcodeByte >> 7
	if fin == 0 {

	}
	log.Println(header.OpcodeByte)
	//读取opcode
	opcode = header.OpcodeByte & 0x0f
	if opcode == 8 {
		//返回error
		return &ProtocalError{
			err_type: Protocal_Error_TYPE_DISCONNECT,
			err:      errors.New("Connection closed"),
		}
	}

	header.PayloadByte, err = buf.ReadByte()

	if err != nil {
		return &ProtocalError{
			err_type: Protocal_Error_TYPE_DISCONNECT,
			err:      err,
		}
	}

	mask := header.PayloadByte >> 7
	payload := header.PayloadByte & 0x7f

	var (
		lengthBytes  []byte
		length       uint64
		l            uint16
		maskKeyBytes []byte
		contentBuf   []byte
	)

	//读取长度
	switch {
	case payload < 126:
		length = uint64(payload)

	case payload == 126:
		lengthBytes = make([]byte, 2)
		buf.Read(lengthBytes)
		binary.Read(bytes.NewReader(lengthBytes), binary.BigEndian, &l)
		length = uint64(l)

	case payload == 127:
		lengthBytes = make([]byte, 8)
		buf.Read(lengthBytes)
		binary.Read(bytes.NewReader(lengthBytes), binary.BigEndian, &length)
	}

	if mask == 1 {
		maskKeyBytes = make([]byte, 4)
		buf.Read(maskKeyBytes)
	}

	contentBuf = make([]byte, length)
	buf.Read(contentBuf)

	if mask == 1 {
		//解码内容
		for i, v := range contentBuf {
			contentBuf[i] = v ^ maskKeyBytes[i%4]
		}
	}
	websocket.Content = contentBuf
	client.callbackServe(SYS_ON_MESSAGE)

	return nil
}

func (this *WebSocketProtocol) Handshake(client *Response) (bool, []byte) {
	if client.conn.HandshakeFlg == true {
		return true, nil
	}
	reader := client.conn.buf
	key := ""
	str := ""
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
			return false, nil
		}
		if len(line) == 0 {
			break
		}
		str = string(line)
		if strings.HasPrefix(str, "Sec-WebSocket-Key") {
			key = str[19:43]
		}
	}
	sha := sha1.New()
	io.WriteString(sha, key+"258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	key = base64.StdEncoding.EncodeToString(sha.Sum(nil))

	header := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Version: 13\r\n" +
		"Sec-WebSocket-Accept: " + key + "\r\n" +
		"Upgrade: websocket\r\n\r\n"
	client.conn.HandshakeFlg = true
	return true, []byte(header)
}

func (gateway *WebSocketProtocol) Encode(msg interface{}) (buf []byte, err error) {
	buff := new(bytes.Buffer)
	//binary.Write(buff, binary.BigEndian, msg.(*GateWayProtocol).Header)

	frame := []byte{129}

	data := msg.(*WebSocketProtocol).Content
	length := len(data)

	switch {
	case length < 126:
		frame = append(frame, byte(length))
	case length <= 0xffff:
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(length))
		frame = append(frame, byte(126))
		frame = append(frame, buf...)
	case uint64(length) <= 0xffffffffffffffff:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(length))
		frame = append(frame, byte(127))
		frame = append(frame, buf...)
	default:
		log.Fatal("Data too large")
		return
	}
	frame = append(frame, data...)

	binary.Write(buff, binary.BigEndian, frame)
	//MyLog().Info(buff.Bytes())
	return buff.Bytes(), nil
}

func (req *WebSocketProtocol) Decode(buf []byte) (err error) {
	return nil
}

func (req *WebSocketProtocol) SetCmd(cmd uint16) {
	//req.Header.Cmd = cmd
}

func (req *WebSocketProtocol) GetCmd() (cmd uint16) {
	return 0
	//return req.Header.Cmd
}
