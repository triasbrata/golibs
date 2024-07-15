package types

import "net"

type ReaderUDP interface {
	IsConClose() bool
	GetCon() *net.UDPConn
	GetMaxLengthMessage() int64
}
type EventInvoke func(payload Dto, remoteAddr net.Addr)
