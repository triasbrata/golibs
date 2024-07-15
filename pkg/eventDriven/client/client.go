package client

import (
	cl "github.com/triasbrata/golibs/pkg/eventDriven/internals/client"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
)

func NewClient() types.Client {
	return cl.NewInternalClient("", nil)
}
