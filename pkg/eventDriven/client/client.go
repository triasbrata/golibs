package client

import (
	cl "github.com/triasbrata/golibs/pkg/eventDriven/internals/client"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
)

type client struct {
	cl.InternalClient
}

func NewClient() types.Client {
	return &client{}
}
