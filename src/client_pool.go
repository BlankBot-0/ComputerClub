package src

type ClientPool struct {
	Pool map[string]Client
}

func NewClientPool() ClientPool {
	return ClientPool{
		Pool: make(map[string]Client),
	}
}

func (cp *ClientPool) GetClient(name string) (Client, bool) {
	client, ok := cp.Pool[name]
	return client, ok
}

func (cp *ClientPool) Add(name string, client Client) error {
	if _, ok := cp.Pool[name]; ok {
		return YouShallNotPass
	}
	cp.Pool[name] = client
	return nil
}

func (cp *ClientPool) Update(name string, client Client) error {
	if _, ok := cp.Pool[name]; !ok {
		return ClientUnknown
	}
	cp.Pool[name] = client
	return nil
}

func (cp *ClientPool) Remove(client string) error {
	if _, ok := cp.Pool[client]; !ok {
		return ClientUnknown
	}
	delete(cp.Pool, client)
	return nil
}

func (cp *ClientPool) ListClients() []string {
	clients := make([]string, 0, len(cp.Pool))
	for client := range cp.Pool {
		clients = append(clients, client)
	}
	return clients
}
