package eventdriven

import (
	"net"
	"testing"

	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
)

func Test_server_reciveMessage(t *testing.T) {
	type fields struct {
		con       *net.UDPConn
		quit      chan struct{}
		close     bool
		events    map[string]map[string]interface{}
		namespace string
		Id        string
		idGen     types.ShortID
		ServOpt   *ServOpt
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				con:       tt.fields.con,
				quit:      tt.fields.quit,
				close:     tt.fields.close,
				events:    tt.fields.events,
				namespace: tt.fields.namespace,
				Id:        tt.fields.Id,
				idGen:     tt.fields.idGen,
				ServOpt:   tt.fields.ServOpt,
			}
			s.receiveMessage()
		})
	}
}
