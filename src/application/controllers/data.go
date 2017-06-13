package controllers

import (
	"application/libraries/helpers"
	"fmt"

	"github.com/beck917/pillX/pillx"
	"github.com/bitly/go-simplejson"
	"github.com/garyburd/redigo/redis"
)

func SubScripTionHandler(client *pillx.Response, protocol pillx.IProtocol) {
	req := protocol.(*pillx.GateWayProtocol)

	//解析content
	_, jsonErr := simplejson.NewJson(req.Content)
	if jsonErr != nil {
		//记录错误
		return
	}

	jsonRet, _ := simplejson.NewJson([]byte(`{}`))
	jsonRet.Set("method", "subscription")
	//jsonRet.Set("match", "valid")
	req.Content, _ = jsonRet.Encode()
	req.Header.Size = uint16(len(req.Content))

	psc := redis.PubSubConn{Conn: helpers.GlobalRedis.Conn}
	psc.Subscribe("brtest")
	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
				jsonRet.Set("data", v.Data)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				fmt.Println(v)
			}
		}
	}()

	//发送到所有网关
	pillx.SendAllGateWay(pillx.GatewayPools, req)
}
