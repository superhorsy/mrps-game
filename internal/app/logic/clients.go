package logic

import (
	"sync"

	"mrps-game/internal/app/logic/model"
)

type Clients struct {
	clients map[uint32]*Client
	mu      sync.RWMutex
}

func NewClients() *Clients {
	return &Clients{clients: make(map[uint32]*Client)}
}

func (c *Clients) Add(client *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clients[client.Id] = client
}

func (c *Clients) Get(clientId uint32) (*Client, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	client, ok := c.clients[clientId]
	return client, ok
}

func (c *Clients) Remove(clientId uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.clients, clientId)
}

func (c *Clients) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.clients)
}

func (c *Clients) GetOpponents(clientId uint32) []model.Opponent {
	c.mu.RLock()
	defer c.mu.RUnlock()
	opponents := make([]model.Opponent, 0, len(c.clients)-1)
	for id, client := range c.clients {
		if id == clientId {
			continue
		}
		opponents = append(opponents, model.Opponent{Id: client.Id, Name: client.Name})
	}
	return opponents
}
