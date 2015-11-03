package pillx

type TextProtocol struct {
	Content			[]byte
}

func (req *TextProtocol) New() (protocol IProtocol) {
	return new(TextProtocol)
}

func (req *TextProtocol) Analyze(client *Response) (err error) {
	if (client.conn.handshake_flg != true) {
		client.conn.handshake_flg = true
		//派发连接通知
		client.callbackServe(SYS_ON_CONNECT)
		return nil
	}
	
	req.Content, _, _ = client.conn.buf.ReadLine()
	
	client.callbackServe(SYS_ON_MESSAGE)
	return nil
}

func (req *TextProtocol) Encode(msg interface{}) (buf []byte, err error) {
	return msg.(*TextProtocol).Content,nil
}

func (req *TextProtocol) Decode(buf []byte) (err error) {
	return nil
}

func (req *TextProtocol) GetCmd() (cmd uint16) {
	return 1
}
