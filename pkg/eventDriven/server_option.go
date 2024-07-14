package eventdriven

import (
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/cons"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
)

type ServOpt struct {
	MaxLenMessage int64
}

// WithMaxLengthMessage implements ServerOptions.
func (so *ServOpt) WithMaxLengthMessage(len int64) {
	so.MaxLenMessage = len
}

func (so *ServOpt) fill() {
	if so.MaxLenMessage == 0 {
		so.MaxLenMessage = 2 * cons.MB
	}

}
func NewOptions() types.ServerOptions {
	return &ServOpt{}
}
