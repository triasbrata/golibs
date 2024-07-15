package types

type ClientManager interface {
	Register(factory ClientManagerFactory) error
	Get(clientID string) (client Client)
}
type ClientManagerFactory func() (cid string, clientInstance Client, err error)
