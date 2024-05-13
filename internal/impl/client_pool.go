package impl

import (
	"yadroTest/internal/event_manager"
)

type ClientMapPool struct {
	Pool map[event_manager.ClientName]event_manager.Client
}

func NewClientPool() *ClientMapPool {
	return &ClientMapPool{
		Pool: make(map[event_manager.ClientName]event_manager.Client),
	}
}

func (cp *ClientMapPool) GetClient(name event_manager.ClientName) (event_manager.Client, bool) {
	client, ok := cp.Pool[name]
	return client, ok
}

func (cp *ClientMapPool) Add(name event_manager.ClientName, client event_manager.Client) error {
	if _, ok := cp.Pool[name]; ok {
		return event_manager.YouShallNotPass
	}
	cp.Pool[name] = client
	return nil
}

func (cp *ClientMapPool) Remove(client event_manager.ClientName) error {
	if _, ok := cp.Pool[client]; !ok {
		return event_manager.ClientUnknown
	}
	delete(cp.Pool, client)
	return nil
}

func (cp *ClientMapPool) ListClients() []event_manager.ClientName {
	clients := make([]event_manager.ClientName, 0, len(cp.Pool))
	for client := range cp.Pool {
		clients = append(clients, client)
	}
	return clients
}
