package pillx

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
	"stathat.com/c/consistent"
)

var globalConsistent *consistent.Consistent = consistent.New()
var logFormat string

func GetPool(workerPools map[string]Pool) (wp Pool, key string) {
	//随机取出一个workerpool
	//实现一个稳定的hash算法
	for key, wp := range workerPools {
		return wp, key
	}
	return
}

func SetLogFormat(format string) {
	logFormat = format
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

var logInitFlg bool = false

func MyLog() *log.Entry {
	if logInitFlg == false {
		// Log as JSON instead of the default ASCII formatter.
		if logFormat != "" {
			log.SetFormatter(&log.JSONFormatter{})
			// Only log the warning severity or above.
			log.SetLevel(log.DebugLevel)
			// Output to stderr instead of stdout, could also be a file.
			file := getFile(logFormat)
			log.SetOutput(file)
		} else {
			// The TextFormatter is default, you don't actually have to do this.
			log.SetFormatter(&log.TextFormatter{})
			// Only log the warning severity or above.
			log.SetLevel(log.DebugLevel)
			// Output to stderr instead of stdout, could also be a file.
			//log.SetOutput(os.Stderr)
		}

		//每天重置下
		if logFormat != "" {
			c := cron.New()
			c.AddFunc("1 0 0 * * *", func() {
				file := getFile(logFormat)
				log.SetOutput(file)
			})
			c.Start()
		}
		logInitFlg = true
	}
	return log.WithFields(log.Fields{
		"prama": "mylog",
	})
}

func getFile(filename string) *os.File {
	var f *os.File
	if checkFileIsExist(filename) { //如果文件存在
		f, _ = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0777) //打开文件
	} else {
		f, _ = os.Create(filename) //创建文件
	}
	return f
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
