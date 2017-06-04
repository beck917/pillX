package pillx

import (
	"fmt"

	"github.com/beck917/pillX/Proto"

	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/clientv3"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

var (
	GatewayPools  map[string]Pool = make(map[string]Pool)
	WorkerClients                 = make(map[uint64]*WorkerClient)
)

type Worker struct {
	InnerAddr   string
	InnerServer *Server
	EtcdClient  *etcd.Client
	WorkerName  string
	WatchName   string
}

type WorkerClient struct {
	IP string
}

func NewWorker(name string) *Worker {
	return &Worker{}
}

func (worker *Worker) innerConnectHandler(response *Response, protocol IProtocol) {
	//log.Info("连接")
}

func (worker *Worker) innerMessageHandler(response *Response, protocol IProtocol) {
}

func (worker *Worker) innerCloseHandler(response *Response, protocol IProtocol) {

}

func (worker *Worker) innerHandShakeHandler(response *Response, protocol IProtocol) {
	req := protocol.(*GateWayProtocol)
	workerHandShark := &Proto.WorkerHandShark{}
	proto.Unmarshal(req.Content, workerHandShark)
	MyLog().WithFields(log.Fields{
		"client_ip": *(workerHandShark.IP),
		"client_id": req.Header.ClientId,
	}).Info("握手成功")
	//记录客户端ip等信息
	WorkerClients[req.Header.ClientId] = &WorkerClient{
		IP: *(workerHandShark.IP),
	}
}

func (worker *Worker) Init() {
	//设置内部通信地址
	worker.InnerServer = &Server{
		Addr:     worker.InnerAddr,
		Protocol: &GateWayProtocol{},
	}
	worker.InnerServer.HandleFunc(SYS_ON_CONNECT, worker.innerConnectHandler)
	worker.InnerServer.HandleFunc(SYS_ON_MESSAGE, worker.innerMessageHandler)
	worker.InnerServer.HandleFunc(SYS_ON_CLOSE, worker.innerCloseHandler)
	worker.InnerServer.HandleFunc(SYS_ON_HANDSHAKE, worker.innerHandShakeHandler)

	go worker.InnerServer.ListenAndServe()
	log.WithFields(log.Fields{
		"addr": worker.InnerAddr,
	}).Info("inner server started")

	//获取gateway地址
	resp, _ := worker.EtcdClient.Get(context.Background(), worker.WatchName, etcd.WithPrefix())
	for _, ev := range resp.Kvs {
		name := string(ev.Value)
		client := &Client{}
		client.Addr = name
		GatewayPools[name], _ = client.Dail()
		log.WithFields(log.Fields{
			"key":  string(ev.Key),
			"name": name,
		}).Info("gateway connected")
	}

	//keepalive
	resp1, err := worker.EtcdClient.Grant(context.TODO(), 3)
	MyLog().Info(resp1)
	if err != nil {
		log.Fatal(err)
	}
	//_, err = worker.EtcdClient.Put(context.TODO(), "foo", "bar", clientv3.WithLease(resp.ID))
	_, err = worker.EtcdClient.Put(context.TODO(), worker.WorkerName, worker.InnerAddr, etcd.WithLease(resp1.ID))
	if err != nil {
		log.Fatal(err)
	}
	// the key 'foo' will be kept forever
	_, kaerr := worker.EtcdClient.KeepAlive(context.TODO(), resp1.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}
	//注册worker
	//worker.EtcdClient.Put(context.Background(), worker.WorkerName, worker.InnerAddr)
	return
}

func (worker *Worker) Watch() {
	//监听gateway
	rch := worker.EtcdClient.Watch(context.Background(), worker.WatchName, etcd.WithPrefix())
	log.Info("etcd watch started")
	for wresp := range rch {
		log.Info("etcd watch started1")
		worker.watchGateways(wresp.Events)
	}
	log.Info("etcd watch started2")
}

func (worker *Worker) watchGateways(events []*etcd.Event) {
	for _, ev := range events {
		name := string(ev.Kv.Value)
		if _, ok := GatewayPools[name]; !ok {
			client := &Client{}
			client.Addr = name
			GatewayPools[name], _ = client.Dail()
		}
		fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
	}
	return
}
