package main

import (
	"Proto"
	"fmt"
	"libraries/toml"
	"libraries/utils"
	"net/http"
	"pillx"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/robfig/cron"
)

type Users struct {
	Msg     string `json:"msg"`
	MsgCode int    `json:"msg_code"`
	Data    struct {
		User struct {
			Black []string `json:"black"`
			Admin []string `json:"admin"`
		} `json:"user"`
	} `json:"data"`
	ServerTime int `json:"server_time"`
	Status     int `json:"status"`
}

func lastMsg(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/x-protobuf")

	//req.ParseForm()

	page := 1
	if req.PostFormValue("page") != "" {
		page, _ = strconv.Atoi(req.PostFormValue("page"))
	}
	pagecount := 30

	msglist := pillx.RecentMsgChan.GetList()

	messageProto := &Proto.ChatList{}
	pagei := 0
	for iter := msglist.Front(); iter != nil; iter = iter.Next() {
		pagei++
		if pagei <= pagecount*(page-1) {
			continue
		}
		if pagei > pagecount*(page) {
			break
		}
		messageProto.ChatItem = append(messageProto.ChatItem, iter.Value.(*Proto.ChatData))
	}
	v, _ := proto.Marshal(messageProto)
	w.Write(v)
}

func getAdminBlack() {
	pillx.BlackIdMap = make(map[int32]bool)
	pillx.AdminIdMap = make(map[int32]bool)
	//存入数据
}

func main() {
	tomlConfig, err := toml.LoadTomlConfig("C:\\Code\\pillX\\bin\\etc\\config.toml")
	if err != nil {
		panic(err)
	}

	etcdClient := utils.EtcdDail(tomlConfig.Etcd)

	http.HandleFunc("/lastmsg", lastMsg)
	go http.ListenAndServe(":8008", nil)

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
		getAdminBlack()
	})
	c.Start()
	getAdminBlack()

	gateway.OuterProtocol = &pillx.WebSocketProtocol{}
	gateway.EtcdClient = etcdClient
	gateway.Init()
}
