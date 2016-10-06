package pillx

import (
	"sync"
	"time"
)

type Channel struct {
	rwmu    sync.RWMutex
	name    string
	clients map[uint64]*Response

	//此频道屏蔽用户 1->频道中的某个客户端id 2->此客户端屏蔽的id字典
	clientBlocks map[uint64]map[uint64]int64

	uidBlocks map[int32]map[int32]int64
}

var pubsubChannels map[string]*Channel
var blocktime int64 = 1800 //屏蔽时间

func NewChannel(name string) *Channel {
	if pubsubChannels == nil {
		pubsubChannels = make(map[string]*Channel)
	}

	if _, ok := pubsubChannels[name]; ok {
		return pubsubChannels[name]
	}

	channel := &Channel{
		name:    name,
		clients: make(map[uint64]*Response),
	}
	channel.clientBlocks = make(map[uint64]map[uint64]int64)
	channel.uidBlocks = make(map[int32]map[int32]int64)
	pubsubChannels[name] = channel
	return channel
}

func (channel *Channel) Block(client *Response, blockClientId uint64) {
	blockStruct, ok := channel.clientBlocks[client.GetConn().Id]
	if !ok {
		blockStruct = make(map[uint64]int64)
		channel.clientBlocks[client.GetConn().Id] = blockStruct
	}
	blockStruct[blockClientId] = time.Now().Unix()
	//MyLog().Info(channel.clientBlocks)
}

func (channel *Channel) BlockUid(uid int32, blockUid int32) {
	blockStruct, ok := channel.uidBlocks[uid]
	if !ok {
		blockStruct = make(map[int32]int64)
		channel.uidBlocks[uid] = blockStruct
	}
	blockStruct[blockUid] = time.Now().Unix()
	MyLog().Info(channel.uidBlocks)
}

func (channel *Channel) Subscribe(client *Response) {
	channel.rwmu.Lock()
	channel.clients[client.conn.Id] = client
	client.channels[channel.name] = channel
	channel.rwmu.Unlock()
}

func (channel *Channel) Publish(client *Response, message interface{}) {
	for _, chanClient := range channel.clients {
		//fmt.Print(client)
		//屏蔽
		//MyLog().Info(channel.clientBlocks[chanClient.GetConn().Id][client.GetConn().Id])
		//MyLog().Info(channel.clientBlocks)

		chanUid, chok := clientMap[chanClient.conn.Id]
		clientUid, _ := clientMap[client.conn.Id]

		var ok bool = false
		var clientBlockTime int64
		if chok {
			clientBlockTime, ok = channel.uidBlocks[chanUid.Uid][clientUid.Uid]
			//MyLog().Info(channel.uidBlocks, chanUid.Uid, clientUid.Uid)
			if ok {
				//判断时间是否过期
				if clientBlockTime+blocktime < time.Now().Unix() {
					//过期，删除
					delete(channel.uidBlocks[chanUid.Uid], clientUid.Uid)
					ok = false
				}
			}
		}
		/**
		clientBlockTime, ok := channel.clientBlocks[chanClient.GetConn().Id][client.GetConn().Id]
		if ok {
			//判断时间是否过期
			if clientBlockTime+blocktime < time.Now().Unix() {
				//过期，删除
				delete(channel.clientBlocks[chanClient.GetConn().Id], client.GetConn().Id)
				ok = false
			}
		}
		*/
		if !ok && client.GetConn().Id != chanClient.GetConn().Id {
			if chanClient != nil {
				//MyLog().Info(chanClient.GetConn().Id, client.GetConn().Id)
				chanClient.Send(message)
			}
		}
	}
}

func (channel *Channel) UnSubscribe(client *Response) {
	channel.rwmu.Lock()
	delete(channel.clients, client.conn.Id)
	delete(channel.clientBlocks, client.conn.Id)
	channel.rwmu.Unlock()
}
