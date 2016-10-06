package main

import (
	"application/controllers"
	"fmt"
	"libraries/opcodes"
	"libraries/toml"
	"libraries/utils"
	"pillx"
)

func main() {
	tomlConfig, err := toml.LoadTomlConfig("etc/config.toml")
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
	worker.Watch()
}
