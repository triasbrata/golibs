package types

type ServerOptions interface {
	WithMaxLengthMessage(len int64)
}
