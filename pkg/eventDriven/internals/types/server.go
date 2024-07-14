package types

import (
	"net"
)

type Server interface {
	Event
	Close()
	Listen(address ...string) error
}
type ServerClientManager interface {
	Register(clientID string, addr *net.UDPAddr) error
	Get(clientID string) (client *Client)
}
