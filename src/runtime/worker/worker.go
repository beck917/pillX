package main

import (
	"fmt"

	"application/controllers"
	"application/libraries/helpers"
	"application/libraries/opcodes"
	"application/libraries/toml"

	"github.com/beck917/pillX/pillx"
	"github.com/garyburd/redigo/redis"
)

func main() {
	tomlConfig, err := toml.LoadTomlConfig("./etc/config.toml")
	if err != nil {
		panic(err)
	}

	psc := redis.PubSubConn{Conn: helpers.GlobalRedis.Conn}
	psc.Subscribe("brtest")
	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				fmt.Println(v)
			}
		}
	}()

	etcdClient := helpers.EtcdDail(tomlConfig.Etcd)

	worker := &pillx.Worker{
		InnerAddr:  fmt.Sprintf("%s:%d", tomlConfig.Pillx.WorkerInnerHost, tomlConfig.Pillx.WorkerInnerPort),
		WorkerName: fmt.Sprintf("%s%d", tomlConfig.Pillx.WorkerName, 1),
		WatchName:  tomlConfig.Pillx.GatewayName,
	}
	helpers.GlobalWorker = worker
	worker.EtcdClient = etcdClient

	worker.Init()
	worker.InnerServer.HandleFunc(opcodes.APP_INDEX, controllers.IndexHandler)
	worker.InnerServer.HandleFunc(opcodes.APP_LOGIN, controllers.LoginHandler)
	worker.InnerServer.HandleFunc(opcodes.APP_BOOKING, controllers.SubScripTionHandler)
	worker.InnerServer.HandleFunc(pillx.SYS_ON_MESSAGE, controllers.OnMessageHandler)
	worker.Watch()
}
