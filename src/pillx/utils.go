package pillx

import log "github.com/Sirupsen/logrus"

func GetPool(workerPools map[string]Pool) (wp Pool, key string) {
	//随机取出一个workerpool
	//实现一个稳定的hash算法
	for key, wp := range workerPools {
		return wp, key
	}
	return
}

func SendAllGateWay(gatewayPools map[string]Pool, msg interface{}) {
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

func MyLog() *log.Entry {
	return log.WithFields(log.Fields{
		"prama": "mylog",
	})
}
