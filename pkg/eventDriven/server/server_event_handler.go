package server

import (
	"fmt"
	"net"

	"github.com/triasbrata/golibs/pkg/eventDriven/internals/client"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/events"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/model"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
)

func handlerConnecting(s *server) types.ClientHandlerDataOnly {
	return func(d any) error {

		if remoteAddr, safe := d.(net.Addr); safe {
			fmt.Printf("here d: %v\n", d)
			clientID := s.idGen.Generate()

			err := s.clients.Register(func() (cid string, cl types.Client, err error) {
				cl = client.NewInternalClient(clientID, remoteAddr)
				return clientID, cl, nil
			})
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return fmt.Errorf("failed register client with error %w", err)
			}
			return s.Send(model.NewDto(remoteAddr, s.Id, clientID, s.namespace, events.CONNECTED, "hello"))
		}
		return fmt.Errorf("failed register client, when try parse remote address")

	}
}
