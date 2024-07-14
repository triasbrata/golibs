package gen

import (
	"github.com/lithammer/shortuuid/v4"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
)

type sid struct {
	alphaNum string
}

// Generate implements types.ShortID.
func (s *sid) Generate() string {
	return shortuuid.NewWithAlphabet(s.alphaNum)
}

func NewSID() types.ShortID {
	return &sid{
		alphaNum: newOrderAlphaNumeric(),
	}
}
