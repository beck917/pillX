package pillx

import (
	"fmt"
)

type Channel struct {
	name 	string
	clients map[*Response] *Response
}

var pubsubChannels map[string] *Channel

func NewChannel(name string) *Channel {
	if pubsubChannels == nil {
		pubsubChannels = make(map[string] *Channel)
	}
	
	if _, ok := pubsubChannels[name]; ok {
		return pubsubChannels[name]
	}
	
	channel := &Channel{
		name: name,
		clients: make(map[*Response]*Response),
	}
	pubsubChannels[name] = channel
	return channel
}

func (channel *Channel) Subscribe (client *Response) {
	channel.clients[client] = client
}

func (channel *Channel) Publish (message *Request) {
	for _, client := range channel.clients {
		fmt.Print(client)
		client.Send(message)
	}
}

func (channel *Channel) UnSubscribe (client *Response) {
	delete(channel.clients, client)
}