package pillx

import (
	"container/list"
	"sync"
)

type MsgList struct {
	mu   sync.RWMutex
	list *list.List
	max  int
}

func NewMsgList(max int) *MsgList {
	msglist := &MsgList{}
	msglist.list = list.New()
	msglist.max = max
	return msglist
}

func (msglist *MsgList) Push(msg interface{}) {
	msglist.mu.Lock()
	defer msglist.mu.Unlock()
	if msglist.list.Len() >= msglist.max {
		msg := msglist.list.Back()
		msglist.list.Remove(msg)
	}
	msglist.list.PushFront(msg)
}

func (msglist *MsgList) GetSlice() (slice []interface{}) {
	for iter := msglist.list.Front(); iter != nil; iter = iter.Next() {
		slice = append(slice, iter.Value)
	}
	return
}

func (msglist *MsgList) GetList() *list.List {
	return msglist.list
}
