package model

import (
	"container/list"
	"fmt"
	"sync"
)

type MessageQueue struct {
	rwMutex *sync.RWMutex
	queue   *list.List
}

var messageQueue *MessageQueue

func NewMessageQueue() *MessageQueue {
	if messageQueue == nil {
		messageQueue = new(MessageQueue)
	}

	return messageQueue
}

func (mq *MessageQueue) Initialize() {
	mq.queue = list.New()
	mq.rwMutex = new(sync.RWMutex)
}

func (mq *MessageQueue) SendMessage(member *Member) {
	mq.rwMutex.Lock()
	defer mq.rwMutex.Unlock()
	mq.queue.PushBack(member)
}

func (mq *MessageQueue) GetMessage() (*Member, error) {
	mq.rwMutex.Lock()
	defer mq.rwMutex.Unlock()
	if mq.queue.Len() > 0 {
		e := mq.queue.Front()
		member := e.Value.(*Member)
		mq.queue.Remove(e)
		return member, nil
	} else {
		return nil, fmt.Errorf("empty queue")
	}
}

func (mq *MessageQueue) Finalize() {
	mq.rwMutex.Lock()
	defer mq.rwMutex.Unlock()

}
