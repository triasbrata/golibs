package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/triasbrata/golibs/pkg/eventDriven/internals/client"
	clientmanager "github.com/triasbrata/golibs/pkg/eventDriven/internals/clientmanagers"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/crypto"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/events"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/model"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/parser"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/validator"
	"github.com/vmihailenco/msgpack/v5"
)

type server struct {
	con   *net.UDPConn
	quit  chan struct{}
	close bool
	// mapping by  event namespace handler
	events map[string]map[string]interface{}
	//holder client
	clients   types.ClientManager
	namespace string
	Id        string
	idGen     types.ShortID
	wg        *sync.WaitGroup
	*ServOpt
}

// GetCon implements types.ReaderUDP.
func (s *server) GetCon() *net.UDPConn {
	return s.con
}

// GetMaxLengthMessage implements types.ReaderUDP.
func (s *server) GetMaxLengthMessage() int64 {
	return s.MaxLenMessage
}

// IsConClose implements types.ReaderUDP.
func (s *server) IsConClose() bool {
	return s.close
}

// Close implements Server.
func (s *server) Close() {
	if !s.close {
		s.wg.Done()
		s.close = true
	}
}

// Event implements Server.
func (s *server) Event(event string, h interface{}) error {
	err := validator.ValidateEvent(h)
	if err != nil {
		return err
	}
	ne, safe := s.events[event]
	if !safe {
		ne = make(map[string]interface{})
	}
	ne[s.namespace] = h
	s.events[event] = ne
	fmt.Printf("event %s registered\n", event)
	return nil
}

func (s *server) findEvent(event string, namespace string) interface{} {
	if ne, safe := s.events[event]; safe {
		h, safe := ne[namespace]
		if safe {
			return h

		}
	}
	return nil
}

// Send implements Server.
func (s *server) Send(data types.Dto) error {
	//send to self
	if data.ReciverID() == data.SenderID() {
		eventHandler := s.findEvent(data.Event(), data.Namespace())
		if eventHandler != nil {

			return client.HandlerInvoker(eventHandler, s.clients.Get(s.Id), data.Data())
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

	_, err = s.con.WriteToUDP(payloads, addr)

	if err != nil {
		return err
	}
	return nil
}

// Listen implements Server.
func (s *server) Listen(address ...string) (err error) {
	serverAddress := ":9040"
	if len(address) == 1 {
		serverAddress = address[0]
	}
	ip, port, err := parser.ParseIpAndPort(serverAddress)
	if err != nil {
		return fmt.Errorf("failed to parse IP and port: %w", err)
	}

	s.con, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   ip,
		Port: port,
	})
	if err != nil {
		return fmt.Errorf("failed to listen on UDP: %w", err)
	}
	defer func() {
		fmt.Printf("server : %v\n", "close")
		s.con.Close()
	}()

	s.close = false
	s.wg.Add(1)
	go s.receiveMessage()
	fmt.Printf("server listening at %v:%v\n", ip, port)

	//self register as client and invoke connected event
	locAddr := s.con.LocalAddr()
	s.clients.Register(func() (cid string, cl types.Client, err error) {
		return s.Id, client.NewInternalClient(s.Id, locAddr), nil
	})
	s.Send(model.NewDto(locAddr, s.Id, s.Id, s.namespace, events.CONNECTED, "hello"))

	//wait all process recive message
	s.wg.Wait()
	return nil
}
func (s *server) receiveMessage() {
	parser.ListenNewMessage(s, func(payload types.Dto, remoteAddr net.Addr) {
		if ne, safe := s.events[payload.Event()]; safe {
			if handler, safe := ne[payload.Namespace()]; safe {
				switch payload.Event() {
				case events.CONNECTING:
					s.clients.Register(func() (cid string, cl types.Client, err error) {
						clientID := s.idGen.Generate()
						return clientID, client.NewInternalClient(clientID, remoteAddr), nil
					})
					client.HandlerInvoker(handler, nil, remoteAddr)
				default:
					client.HandlerInvoker(handler, s.clients.Get(payload.SenderID()), payload.Event())
				}
			}
		}
	})
}

func New(option types.ServerOptions) (types.Server, error) {
	var so types.ServerOptions = &ServOpt{}
	if option != nil {
		so = option
	}
	servOpt := so.(*ServOpt)
	servOpt.fill()
	sid := crypto.NewSID()

	serverInstance := &server{
		quit:      make(chan struct{}, 1),
		ServOpt:   servOpt,
		close:     true,
		clients:   clientmanager.NewMemoryManager(),
		events:    make(map[string]map[string]interface{}),
		Id:        sid.Generate(),
		idGen:     sid,
		namespace: "/",
		wg:        &sync.WaitGroup{},
	}
	serverInstance.hookInit()

	return serverInstance, nil
}
func (s *server) hookInit() {
	s.Event(events.CONNECTING, handlerConnecting(s))
}
