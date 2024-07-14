package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	eventdriven "github.com/triasbrata/golibs/pkg/eventDriven"
	"github.com/triasbrata/golibs/pkg/eventDriven/client"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/events"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
)

func Test_initServer(t *testing.T) {
	e, _ := eventdriven.New(eventdriven.NewOptions())
	err := e.Event(events.CONNECTED, types.ClientHandler(func(sender types.Client, message any) error {
		msg, safe := message.(string)
		assert.Equal(t, safe, true)
		assert.Equal(t, msg, "hello")
		assert.Equal(t, fmt.Sprintf("%T", sender), "*client.InternalClient")
		e.Close()
		return nil
	}))
	assert.Nil(t, err)
	err = e.Listen()
	assert.Nil(t, err)

}
func Test_initClientAfterServer(t *testing.T) {
	e, _ := eventdriven.New(eventdriven.NewOptions())

	err := e.Event(events.CONNECTED, func(sender types.Client, message any) error {
		c := client.NewClient()
		err := c.Open(":9040")
		assert.Nil(t, err)
		e.Close()
		return nil
	})
	assert.Nil(t, err)
	err = e.Listen("")
	assert.Nil(t, err)

}
