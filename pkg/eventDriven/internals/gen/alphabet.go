package gen

import (
	"math/rand"

	"github.com/triasbrata/golibs/pkg/eventDriven/internals/cons"
)

func newOrderAlphaNumeric() string {
	str := []byte(cons.ALPHA_NUMERIC)
	rand.Shuffle(len(cons.ALPHA_NUMERIC), func(i, j int) {
		temp := str[i]
		strj := str[j]
		str[j] = temp
		str[i] = strj
	})
	return string(str)
}
