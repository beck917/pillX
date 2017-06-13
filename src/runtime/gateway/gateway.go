package main

import (
	"fmt"

	"application/libraries/helpers"
	"application/libraries/toml"

	"github.com/beck917/pillX/pillx"
	"github.com/robfig/cron"
)

func main() {
	tomlConfig, err := toml.LoadTomlConfig("./etc/config.toml")
	if err != nil {
		panic(err)
	}

	etcdClient := helpers.EtcdDail(tomlConfig.Etcd)

	gateway := &pillx.GatewayWebsocket{
		InnerAddr:   fmt.Sprintf("%s:%d", tomlConfig.Pillx.GatewayInnerHost, tomlConfig.Pillx.GatewayInnerPort),
		OuterAddr:   fmt.Sprintf("%s:%d", tomlConfig.Pillx.GatewayOuterHost, tomlConfig.Pillx.GatewayOuterPort),
		GatewayName: fmt.Sprintf("%s%d", tomlConfig.Pillx.GatewayName, 1),
		WatchName:   tomlConfig.Pillx.WorkerName,
	}

	//获取管理员和
	c := cron.New()
	c.AddFunc("0 */5 * * * *", func() {
		pillx.MyLog().Info("获取管理员和黑名单id ")
	})
	c.Start()

	gateway.OuterProtocol = &pillx.WebSocketProtocol{}
	gateway.EtcdClient = etcdClient
	gateway.Init()
}
