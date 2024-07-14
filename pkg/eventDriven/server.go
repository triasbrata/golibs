package eventdriven

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/triasbrata/golibs/pkg/eventDriven/internals/client"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/events"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/gen"
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
	clients   types.ServerClientManager
	namespace string
	Id        string
	idGen     types.ShortID
	wg        *sync.WaitGroup
	*ServOpt
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
	fmt.Printf("event %s registered", event)
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

	if data.ReciverID() == s.Id {
		eventHandler := s.findEvent(data.Event(), data.Namespace())
		fmt.Printf("eventHandler: %v\n", eventHandler)
		if eventHandler != nil {
			return client.HandlerInvoker(eventHandler, client.NewInternalClient(s.Id), data.Data())
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

	if addr, ok := s.con.LocalAddr().(*net.UDPAddr); ok {
		s.Send(model.NewDto(*addr, s.Id, s.Id, s.namespace, events.CONNECTED, "hello"))
	}
	s.wg.Wait()
	return nil
}
func (s *server) receiveMessage() {
	for {
		if s.close {
			return
		}
		remoteAddr, err, msg := s.readMessage()
		if err != nil {
			if strings.Contains(err.Error(), "closed network connection") && s.close {
				return
			}
			fmt.Printf("Some error  %v", err)
			return
		}
		payload := model.NewDto(net.UDPAddr{}, "", "", "", "", nil)
		err = msgpack.Unmarshal(msg, payload)
		if err != nil {
			fmt.Printf("Some error when parse %v", err)

		}
		if ne, safe := s.events[payload.Event()]; safe {
			if handler, safe := ne[payload.Namespace()]; safe {
				switch payload.Event() {
				case events.CONNECTING:
					client.HandlerInvoker(handler, nil, remoteAddr)
				default:
					client.HandlerInvoker(handler, client.NewInternalClient(payload.SenderID()), payload.Event())
				}

			}
		}
	}
}

func (s *server) readMessage() (*net.UDPAddr, error, []byte) {
	msg := make([]byte, s.ServOpt.MaxLenMessage)
	_, remoteAddr, err := s.con.ReadFromUDP(msg)
	trimMsg := make([]byte, 0)
	for _, b := range msg {
		if b != 0 {
			trimMsg = append(trimMsg, b)
		}
	}
	return remoteAddr, err, trimMsg
}

func New(option types.ServerOptions) (types.Server, error) {
	var so types.ServerOptions = &ServOpt{}
	if option != nil {
		so = option
	}
	servOpt := so.(*ServOpt)
	servOpt.fill()
	sid := gen.NewSID()

	serverInstance := &server{
		quit:      make(chan struct{}, 1),
		ServOpt:   servOpt,
		close:     true,
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
	var hConnecting types.ClientHandlerDataOnly = func(d any) error {
		if remoteAddr, safe := d.(*net.UDPAddr); safe {
			err := s.clients.Register(s.idGen.Generate(), remoteAddr)
			if err != nil {
				return fmt.Errorf("failed register client with error %w", err)
			}
		}
		return nil
	}
	s.Event(events.CONNECTING, hConnecting)
}
