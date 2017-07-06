package pillx

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	"stathat.com/c/consistent"
)

var globalConsistent *consistent.Consistent = consistent.New()

func GetPool(workerPools map[string]Pool) (wp Pool, key string) {
	//随机取出一个workerpool
	//实现一个稳定的hash算法
	for key, wp := range workerPools {
		return wp, key
	}
	return
}

func responseSend(clientId uint64, resMap map[string]*Response, msg interface{}) (n int, err error) {
	//TODO 发送中发现连接不可用,则剔除

	_, res := GetResponse(clientId, resMap)
	n, err = res.Send(msg)
	return
	/**
	//循环次数
	for i := 0; i < 5; i++ {
		MyLog().Info("teee")
		_, res := GetResponse(resMap)
		MyLog().Info(resMap)
		n, err = res.Send(msg)

		if err != nil {
			MyLog().Error(err)
			continue
		}
		return n, err
	}
	return 0, err
	*/
}

func GetResponse(clientId uint64, resMap map[string]*Response) (ip string, res *Response) {
	//取得一个稳定的节点
	clientIdStr := strconv.Itoa(int(clientId))
	serverName, _ := globalConsistent.Get(clientIdStr)
	return serverName, resMap[serverName]
	/**
	//实现一个稳定的LoadBalance算法
	//权重随机
	for ip, res := range resMap {
		return ip, res
	}
	return
	*/
}

func SendAllGateWay(resMap map[string]*Response, msg interface{}) {
	for _, response := range resMap {
		response.Send(msg)
	}
	return
}

func SendAllGateWayPool(gatewayPools map[string]Pool, msg interface{}) {
	for _, gp := range gatewayPools {
		gateway, err := gp.Get()
		if err != nil {
			log.WithError(err).Error("gateway池返回错误")
		}
		gateway.response.Send(msg)
		//回收
		gateway.Close()
	}
	return
}

func NewGatewayClient(addr string) *Server {
	client := &Server{
		Addr:     addr,
		Handler:  NewServeRouter(),
		Protocol: &GateWayProtocol{},
	}
	return client
}

func MyLog() *log.Entry {
	return log.WithFields(log.Fields{
		"prama": "mylog",
	})
}
