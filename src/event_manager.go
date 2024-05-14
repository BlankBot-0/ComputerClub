package src

import (
	"time"
)

type EventManager struct {
	ClientPool  ClientPool
	DeskStorage *DeskStorage
	ClientQueue ClientQueue
}

func (m *EventManager) ClientArrived(name string) error {
	_, ok := m.ClientPool.GetClient(name)
	if ok {
		return YouShallNotPass
	}
	return m.ClientPool.Add(name, Client{
		Desk:    0,
		AtQueue: false,
	})
}

func (m *EventManager) ClientSatAtTheDesk(deskNum int, name string, currentTime time.Time) error {
	_, ok := m.ClientPool.GetClient(name)
	if !ok {
		return ClientUnknown
	}
	err := m.DeskStorage.Occupy(name, deskNum, currentTime)
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
func (m *EventManager) ClientAwaits(name string) error {
	client, ok := m.ClientPool.GetClient(name)
	if !ok {
		return ClientUnknown
	}
	if client.AtQueue {
		return nil
	}
	if _, ok := m.DeskStorage.FindAvailable(); ok {
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

func (m *EventManager) ClientLeaves(name string, currentTime time.Time) error {
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
		if err := m.DeskStorage.Free(name, currentTime); err != nil {
			return err
		}
	}
	err := m.ClientPool.Remove(name)
	if err != nil {
		return err
	}
	return nil
}
