package impl

import (
	"container/list"
	"yadroTest/internal/event_manager"
)

type ClientChanQueue struct {
	Queue chan event_manager.ClientName
}

func NewClientQueue(n int) *ClientChanQueue {
	return &ClientChanQueue{
		Queue: make(chan event_manager.ClientName, n),
	}
}

func (q *ClientChanQueue) Enqueue(name event_manager.ClientName) error {
	select {
	case q.Queue <- name:
		return nil
	default:
		return event_manager.QueueIsFull
	}
}

func (q *ClientChanQueue) Dequeue() (event_manager.ClientName, error) {
	select {
	case name := <-q.Queue:
		return name, nil
	default:
		return "", event_manager.QueueIsEmpty
	}
}

type ClientListQueue struct {
	Queue     *list.List
	QueueSize int
}

func NewClientListQueue() *ClientListQueue {
	return &ClientListQueue{
		Queue: list.New(),
	}
}

func (q *ClientListQueue) Enqueue(name event_manager.ClientName) error {
	if q.Queue.Len() >= q.QueueSize {
		return event_manager.QueueIsFull
	}
	q.Queue.PushBack(name)
	return nil
}

func (q *ClientListQueue) Dequeue() (event_manager.ClientName, error) {
	name := q.Queue.Front().Value.(event_manager.ClientName)
	q.Queue.Remove(q.Queue.Front())
	return name, nil
}

func (q *ClientListQueue) Remove(name event_manager.ClientName) error {
	element := q.Queue.Front()
	for element != nil && element.Value.(event_manager.ClientName) != name {
		element = element.Next()
	}
	if element == nil {
		return event_manager.ClientNotInQueue
	}
	q.Queue.Remove(element)
	return nil
}
