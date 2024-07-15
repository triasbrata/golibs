package model

import (
	"net"

	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
)

type StdDto struct {
	FAddress   net.Addr    `msgpack:"target"`
	FSenderId  string      `msgpack:"sender_id"`
	FReciverId string      `msgpack:"reciver_id"`
	FEvent     string      `msgpack:"event"`
	FData      interface{} `msgpack:"data",omitempty`
	FNamespace string      `msgpack:"namespace"`
}

// Data implements types.Dto.
func (sd *StdDto) Data() interface{} {
	return sd.FData
}

// Event implements types.Dto.
func (sd *StdDto) Event() string {
	return sd.FEvent
}

// Namespace implements types.Dto.
func (sd *StdDto) Namespace() string {
	return sd.FNamespace
}

// ReciverID implements types.Dto.
func (sd *StdDto) ReciverID() string {
	return sd.FReciverId
}

// SenderID implements types.Dto.
func (sd *StdDto) SenderID() string {
	return sd.FSenderId
}

func (sd *StdDto) Address() net.Addr {
	return sd.FAddress
}

func NewDto(
	address net.Addr,
	sender_id string,
	reciver_id string,
	namespace string,
	event string,
	data interface{},
) types.Dto {
	return &StdDto{
		FAddress:   address,
		FSenderId:  sender_id,
		FReciverId: reciver_id,
		FEvent:     event,
		FData:      data,
		FNamespace: namespace,
	}
}
