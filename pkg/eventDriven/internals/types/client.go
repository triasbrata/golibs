package types

import "net"

type Client interface {
	Event
	Open(serverAddress string) error
	Close() error
}
type Event interface {
	Send(data Dto) error
	Event(event string, h interface{}) error
}
type Dto interface {
	Address() net.Addr
	SenderID() string
	ReciverID() string
	Event() string
	Data() interface{}
	Namespace() string
}

type ClientHandlerDataOnly = func(d any) error
type ClientHandler = func(client Client, d any) error
type ClientHandlerNoParams = func() error
