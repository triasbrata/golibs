package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triasbrata/golibs/pkg/eventDriven/client"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/events"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/model"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
	"github.com/triasbrata/golibs/pkg/eventDriven/server"
)

func Test_initServer(t *testing.T) {
	e, _ := server.New(server.NewOptions())
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
	e, _ := server.New(server.NewOptions())
	c := client.NewClient()
	err := c.Event(events.CONNECTED, types.ClientHandler(func(cl types.Client, d any) error {
		defer func() {
			e.Close()
		}()
		assert.Equal(t, e.ID(), cl.ID())
		return nil
	}))
	assert.Nil(t, err)

	err = e.Event(events.CONNECTED, func(sender types.Client, message any) error {
		err = c.Open("127.0.0.1:9040")
		assert.Nil(t, err)
		return nil
	})
	assert.Nil(t, err)
	err = e.Listen("")
	assert.Nil(t, err)

}
func Test_pingPong(t *testing.T) {
	e, _ := server.New(server.NewOptions())
	c := client.NewClient()
	err := c.Event(events.CONNECTED, types.ClientHandler(func(cl types.Client, d any) error {
		assert.Equal(t, e.ID(), cl.ID())
		cl.Send()
		return nil
	}))
	assert.Nil(t, err)

	err = e.Event(events.CONNECTED, func(sender types.Client, message any) error {
		err = c.Open("127.0.0.1:9040")
		assert.Nil(t, err)
		return nil
	})
	assert.Nil(t, err)
	err = e.Listen("")
	assert.Nil(t, err)

}
