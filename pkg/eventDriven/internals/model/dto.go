package model

import (
	"net"

	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
)

type stdDto struct {
	address    net.UDPAddr `msgpack:"target"`
	sender_id  string      `msgpack:"sender_id"`
	reciver_id string      `msgpack:"reciver_id"`
	event      string      `msgpack:"event"`
	data       interface{} `msgpack:"data"`
	namespace  string
}

// Data implements types.Dto.
func (sd *stdDto) Data() interface{} {
	return sd.data
}

// Event implements types.Dto.
func (sd *stdDto) Event() string {
	return sd.event
}

// Namespace implements types.Dto.
func (sd *stdDto) Namespace() string {
	return sd.namespace
}

// ReciverID implements types.Dto.
func (sd *stdDto) ReciverID() string {
	return sd.reciver_id
}

// SenderID implements types.Dto.
func (sd *stdDto) SenderID() string {
	return sd.sender_id
}

func (sd *stdDto) Address() net.Addr {
	return &sd.address
}

func NewDto(
	address net.UDPAddr,
	sender_id string,
	reciver_id string,
	namespace string,
	event string,
	data interface{},
) types.Dto {
	return &stdDto{
		address:    address,
		sender_id:  sender_id,
		reciver_id: reciver_id,
		event:      event,
		data:       data,
		namespace:  namespace,
	}
}
