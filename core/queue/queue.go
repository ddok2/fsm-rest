package queue

import (
	"container/list"
	"fmt"
	"sync"
)

// MessageQueue represents an queue.
type MessageQueue struct {
	rwMutex *sync.RWMutex
	queue   *list.List
}

var messageQueue *MessageQueue

// NewMessageQueue return an initialized empty queue.
func NewMessageQueue() *MessageQueue {
	if messageQueue == nil {
		messageQueue = &MessageQueue{
			queue:   list.New(),
			rwMutex: new(sync.RWMutex),
		}
	}
	return messageQueue
}

// Initialize initializes queue.
func (mq *MessageQueue) Initialize() {
	mq.queue = list.New()
	mq.rwMutex = new(sync.RWMutex)
}

// Put inserts a new value at the back of queue.
func (mq *MessageQueue) Put(item interface{}) {
	mq.rwMutex.Lock()
	defer mq.rwMutex.Unlock()
	mq.queue.PushBack(item)
}

// Get removes and returns the first element of queue or nil / error
func (mq *MessageQueue) Get() (interface{}, error) {
	mq.rwMutex.Lock()
	defer mq.rwMutex.Unlock()
	if mq.queue.Len() < 1 {
		return nil, fmt.Errorf("empty queue")
	}
	e := mq.queue.Front()
	item := e.Value
	mq.queue.Remove(e)
	return item, nil
}

// Size returns the number of elements of queue.
func (mq *MessageQueue) Size() int {
	mq.rwMutex.Lock()
	defer mq.rwMutex.Unlock()
	return mq.queue.Len()
}
