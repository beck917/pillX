package utils

import (
	"container/list"
	"reflect"
	"sync"
)

type Queue struct {
	mu        sync.RWMutex
	queueList *list.List
}

func NewQueue() *Queue {
	this := &Queue{}
	this.queueList = list.New()

	return this
}

func (this *Queue) Push(msg interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.queueList.PushBack(msg)
}

func (this *Queue) PushList(msgList interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	s := reflect.ValueOf(msgList)

	for i := 0; i < s.Len(); i++ {
		this.queueList.PushBack(s.Index(i).Interface())
	}
}

func (this *Queue) Pop() (v interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()
	element := this.queueList.Front()

	if element != nil {
		v = element.Value
		//删除
		this.queueList.Remove(element)
	} else {
		v = nil
	}
	return
}

func (this *Queue) GetSliceLimit(limit int) (slice []interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	var next *list.Element
	for iter := this.queueList.Front(); iter != nil; iter = next {
		if limit <= 0 {
			return
		}
		slice = append(slice, iter.Value)
		next = iter.Next()
		this.queueList.Remove(iter)
		limit--
	}

	return
}

func (this *Queue) GetList() *list.List {
	return this.queueList
}
