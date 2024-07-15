package types

type Server interface {
	Event
	Close()
	Listen(address ...string) error
}
