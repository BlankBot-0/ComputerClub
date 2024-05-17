package src

import (
	"container/list"
)

type ClientQueue struct {
	Queue     *list.List
	QueueSize int
}

func NewClientListQueue() ClientQueue {
	return ClientQueue{
		Queue: list.New(),
	}
}

func (q *ClientQueue) Enqueue(name string) error {
	if q.Queue.Len() > q.QueueSize {
		return QueueIsFull
	}
	q.Queue.PushBack(name)
	return nil
}

func (q *ClientQueue) Dequeue() (string, bool) {
	name := q.Queue.Front()
	if name == nil {
		return "", false
	}
	q.Queue.Remove(q.Queue.Front())
	return name.Value.(string), true
}

func (q *ClientQueue) Remove(name string) error {
	element := q.Queue.Front()
	for element != nil && element.Value.(string) != name {
		element = element.Next()
	}
	if element == nil {
		return ClientNotInQueue
	}
	q.Queue.Remove(element)
	return nil
}
