package clientmanager

import (
	"fmt"

	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
)

type memoryClientManager struct {
	clients map[string]types.Client
}

// Get implements types.ServerClientManager.
func (m *memoryClientManager) Get(clientID string) (client types.Client) {
	if client, safe := m.clients[clientID]; safe {
		return client
	}
	return nil
}

// Register implements types.ServerClientManager.
func (m *memoryClientManager) Register(factory types.ClientManagerFactory) error {
	clientID, cl, err := factory()
	if err != nil {
		return fmt.Errorf("fail when register client with error from factory: %w", err)
	}
	m.clients[clientID] = cl
	return nil
}

func NewMemoryManager() types.ClientManager {
	return &memoryClientManager{
		clients: make(map[string]types.Client),
	}
}
