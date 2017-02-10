package pillx

import (
	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/clientv3"
)

type GatewayWebsocket struct {
	InnerAddr     string
	OuterAddr     string
	OuterProtocol IProtocol
	InnerServer   *Server
	OuterServer   *Server
	EtcdClient    *etcd.Client
	GatewayName   string
	WatchName     string
}

func (gateway *GatewayWebsocket) outerConnectHandler(client *Response, protocol IProtocol) {
	clients[client.GetConn().Id] = client
	MyLog().WithFields(log.Fields{
		"client_id": client.GetConn().Id,
		"client_ip": client.GetConn().remonte_conn.RemoteAddr(),
	}).Info("连接到网关")
}

func (gateway *GatewayWebsocket) outerMessageHandler(client *Response, protocol IProtocol) {
	req := protocol.(*WebSocketProtocol)

	MyLog().WithFields(log.Fields{
		"client_id": client.GetConn().Id,
		"content":   string(req.Content),
		"client_ip": client.GetConn().remonte_conn.RemoteAddr(),
	}).Info("发送给worker")

	req.Content = []byte("test")
	client.Send(req)
}

func (gateway *GatewayWebsocket) outerCloseHandler(client *Response, protocol IProtocol) {

	//日志
	MyLog().WithFields(log.Fields{
		"client_id": client.GetConn().Id,
		"client_ip": client.GetConn().remonte_conn.RemoteAddr(),
		"info_code": "close1",
	}).Info("连接断开")
}

func (gateway *GatewayWebsocket) Init() {
	MyLog().WithFields(log.Fields{
		"addr": gateway.InnerAddr,
	}).Info("inner server started")

	gateway.OuterServer = &Server{
		Addr:     gateway.OuterAddr,
		Handler:  nil,
		Protocol: gateway.OuterProtocol,
	}
	gateway.OuterServer.HandleFunc(SYS_ON_CONNECT, gateway.outerConnectHandler)
	gateway.OuterServer.HandleFunc(SYS_ON_MESSAGE, gateway.outerMessageHandler)
	gateway.OuterServer.HandleFunc(SYS_ON_CLOSE, gateway.outerCloseHandler)
	//gateway.OuterServer.HandleFunc(SYS_ON_HANDSHAKE, gateway.outerHandShakeHandler)
	gateway.OuterServer.ListenAndServe()
}
