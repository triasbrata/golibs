package types

type Server interface {
	Event
	Close()
	ID() string
	Listen(address ...string) error
}
