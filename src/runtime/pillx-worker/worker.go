package main

import (
	"application/controllers"
	"fmt"
	"github.com/beck917/pillX/libraries/opcodes"
	"github.com/beck917/pillX/libraries/toml"
	"github.com/beck917/pillX/libraries/utils"
	"github.com/beck917/pillX/pillx"
)

func main() {
	tomlConfig, err := toml.LoadTomlConfig("../etc/config.toml")
	if err != nil {
		panic(err)
	}

	etcdClient := utils.EtcdDail(tomlConfig.Etcd)

	worker := &pillx.Worker{
		InnerAddr:  fmt.Sprintf("%s:%d", tomlConfig.Pillx.WorkerInnerHost, tomlConfig.Pillx.WorkerInnerPort),
		WorkerName: fmt.Sprintf("%s%d", tomlConfig.Pillx.WorkerName, 1),
		WatchName:  tomlConfig.Pillx.GatewayName,
	}
	worker.EtcdClient = etcdClient

	worker.Init()
	worker.InnerServer.HandleFunc(opcodes.APP_INDEX, controllers.IndexHandler)
	worker.InnerServer.HandleFunc(opcodes.APP_LOGIN, controllers.LoginHandler)
	worker.InnerServer.HandleFunc(opcodes.APP_BOOKING, controllers.LoginHandler)
	worker.Watch()
}
