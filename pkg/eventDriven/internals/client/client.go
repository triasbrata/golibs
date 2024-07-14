package client

import (
	"fmt"
	"net"

	"github.com/triasbrata/golibs/pkg/eventDriven/internals/events"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/model"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/parser"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
	"github.com/vmihailenco/msgpack/v5"
)

type InternalClient struct {
	Id        string
	connected bool
	con       *net.UDPConn
	namespace string
	rcAdress  *net.UDPAddr
	events    map[string]map[string]interface{}
}

// Close implements types.InternalClient.
func (c *InternalClient) Close() error {
	panic("unimplemented")
}

// Event implements types.InternalClient.
func (c *InternalClient) Event(event string, h interface{}) error {
	panic("unimplemented")
}

// Open implements types.InternalClient.
func (c *InternalClient) Open(serverAddress string) (err error) {
	ip, port, err := parser.ParseIpAndPort(serverAddress)
	if err != nil {
		return fmt.Errorf("failed to parse IP and port: %w", err)
	}

	remoteAddr := &net.UDPAddr{
		IP:   ip,
		Port: port,
	}

	c.con, err = net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		return fmt.Errorf("failed to dial UDP: %w", err)
	}
	return c.Send(model.NewDto(
		*remoteAddr, "", "", c.namespace, events.CONNECTING, nil,
	))
}
func (c *InternalClient) getAdress() *net.UDPAddr {
	if c.rcAdress != nil {
		return c.rcAdress
	}
	var err error
	c.rcAdress, err = net.ResolveUDPAddr(c.con.LocalAddr().Network(), c.con.LocalAddr().String())
	if err != nil {
		panic(err)
	}
	return c.rcAdress
}

func (s *InternalClient) findEvent(event string, namespace string) interface{} {
	if ne, safe := s.events[event]; safe {
		h, safe := ne[namespace]
		if safe {
			return h

		}
	}
	return nil
}

// Send implements types.InternalClient.
func (c *InternalClient) Send(data types.Dto) error {
	if data.ReciverID() == c.Id {
		eventHandler := c.findEvent(data.Event(), data.Namespace())
		if eventHandler != nil {
			return HandlerInvoker(eventHandler, c, data.Data())
		}
		return nil
	}
	payloads, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}
	addr, err := net.ResolveUDPAddr(data.Address().Network(), data.Address().String())
	if err != nil {
		return err
	}

	_, err = c.con.WriteToUDP(payloads, addr)

	if err != nil {
		return err
	}
	return nil
}

func NewInternalClient(id string) types.Client {
	return &InternalClient{
		Id:     id,
		events: make(map[string]map[string]interface{}),
	}
}
