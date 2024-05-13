package event_manager

import (
	"time"
)

type (
	ClientQueue interface {
		Enqueue(name ClientName) error
		Dequeue() (ClientName, error)
		Remove(name ClientName) error
	}
	ClientPool interface {
		GetClient(name ClientName) (Client, bool)
		Add(name ClientName, client Client) error
		Remove(client ClientName) error
		ListClients() []ClientName
	}
	DeskManager interface {
		Occupy(deskNum int, currentTime time.Time) error
		Free(deskNum int, currentTime time.Time) error
		FindAvailable() (int, error)
		AreAnyFree() bool
	}
)

type Deps struct {
	ClientPool  ClientPool
	DeskManager DeskManager
	ClientQueue ClientQueue
}

type EventManager struct {
	Deps
}

func NewEventManager(deps Deps) *EventManager {
	return &EventManager{
		Deps: deps,
	}
}

func (m *EventManager) ClientArrived(name ClientName) error {
	_, ok := m.ClientPool.GetClient(name)
	if ok {
		return YouShallNotPass
	}
	return m.ClientPool.Add(name, Client{
		Desk:    0,
		AtQueue: false,
	})
}

func (m *EventManager) ClientSatAtTheDesk(deskNum int, name ClientName) error {
	_, ok := m.ClientPool.GetClient(name)
	if !ok {
		return ClientUnknown
	}
	err := m.DeskManager.Occupy(deskNum, time.Now())
	if err != nil {
		return err
	}
	err = m.ClientPool.Add(name, Client{
		Desk:    deskNum,
		AtQueue: false,
	})
	return err
}

// ClientAwaits does nothing if the client is already at the queue
func (m *EventManager) ClientAwaits(name ClientName) error {
	client, ok := m.ClientPool.GetClient(name)
	if !ok {
		return ClientUnknown
	}
	if client.AtQueue {
		return nil
	}
	if m.DeskManager.AreAnyFree() {
		return ICanWaitNoLonger
	}

	err := m.ClientQueue.Enqueue(name)
	if err != nil {
		return err
	}
	err = m.ClientPool.Add(name, Client{
		Desk:    0,
		AtQueue: true,
	})
	return err
}

func (m *EventManager) ClientLeaves(name ClientName, currentTime time.Time) error {
	client, ok := m.ClientPool.GetClient(name)
	if !ok {
		return ClientUnknown
	}
	if client.AtQueue {
		if err := m.ClientQueue.Remove(name); err != nil {
			return err
		}
	}
	if client.Desk != -1 {
		if err := m.DeskManager.Free(client.Desk, currentTime); err != nil {
			return err
		}
	}
	err := m.ClientPool.Remove(name)
	if err != nil {
		return err
	}
	return nil
}
