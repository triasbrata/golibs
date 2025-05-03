package client

import (
	"fmt"
	"net"

	clientmanager "github.com/triasbrata/golibs/pkg/eventDriven/internals/clientmanagers"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/cons"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/events"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/model"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/parser"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/validator"
	"github.com/vmihailenco/msgpack/v5"
)

type InternalClient struct {
	Id               string
	connected        bool
	con              *net.UDPConn
	namespace        string
	rcAdress         net.Addr
	clients          types.ClientManager
	events           map[string]map[string]interface{}
	maxLengthMessage int64
	closed           bool
}

// Send implements types.Client.
func (c *InternalClient) Send(event string, data interface{}) error {
	return c.CommitMsg(&model.StdDto{
		FAddress: ,
	})
}

// ID implements types.Client.
func (c *InternalClient) ID() string {
	return c.Id
}

// GetCon implements types.ReaderUDP.
func (c *InternalClient) GetCon() *net.UDPConn {
	return c.con
}

// GetMaxLengthMessage implements types.ReaderUDP.
func (c *InternalClient) GetMaxLengthMessage() int64 {
	return c.maxLengthMessage
}

// IsConClose implements types.ReaderUDP.
func (c *InternalClient) IsConClose() bool {
	return c.closed
}

// Close implements types.InternalClient.
func (c *InternalClient) Close() error {
	panic("unimplemented")
}

// Event implements types.InternalClient.
func (c *InternalClient) Event(event string, h interface{}) error {
	err := validator.ValidateEvent(h)
	if err != nil {
		return err
	}
	ne, safe := c.events[event]
	if !safe {
		ne = make(map[string]interface{})
	}
	fmt.Printf("c.events: %v\n", c.events[event])
	fmt.Printf("ne: %v\n", ne)
	ne[c.namespace] = h
	c.events[event] = ne
	return nil

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
	c.closed = false
	go c.fetchMessage(c.con.LocalAddr())
	return c.CommitMsg(model.NewDto(
		remoteAddr, "", "", c.namespace, events.CONNECTING, nil,
	))
}
func (c *InternalClient) fetchMessage(localAddr net.Addr) {
	parser.ListenNewMessage(c, func(payload types.Dto, remoteAddr net.Addr) {
		if ne, safe := c.events[payload.Event()]; safe {
			if handler, safe := ne[payload.Namespace()]; safe {
				switch payload.Event() {
				case events.CONNECTED:
					//register host
					c.clients.Register(func() (cid string, clientInstance types.Client, err error) {
						return payload.SenderID(), NewInternalClient(payload.SenderID(), remoteAddr), nil
					})
					//register self
					c.clients.Register(func() (cid string, clientInstance types.Client, err error) {
						return payload.ReciverID(), NewInternalClient(payload.ReciverID(), localAddr), nil
					})
					HandlerInvoker(handler, c.clients.Get(payload.SenderID()), payload.Event())
				default:
					HandlerInvoker(handler, c.clients.Get(payload.SenderID()), payload.Event())
				}
			}
		}
	})
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
func (c *InternalClient) CommitMsg(data types.Dto) error {
	if data.ReciverID() == data.SenderID() && data.ReciverID() != "" {
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
	_, err = c.con.Write(payloads)
	if err != nil {
		return err
	}
	return nil
}

func NewInternalClient(id string, address net.Addr) types.Client {

	return &InternalClient{
		Id:               id,
		events:           make(map[string]map[string]interface{}),
		rcAdress:         address,
		namespace:        "/",
		clients:          clientmanager.NewMemoryManager(),
		connected:        false,
		con:              nil,
		maxLengthMessage: 5 * cons.MB,
		closed:           true,
	}
}
