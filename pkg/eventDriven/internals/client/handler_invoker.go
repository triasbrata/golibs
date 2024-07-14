package client

import (
	"fmt"
	"reflect"

	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
)

func HandlerInvoker(h interface{}, client types.Client, value interface{}) error {
	vh := reflect.ValueOf(h)
	if vh.Kind() == reflect.Func {
		in := []reflect.Value{}
		numParam := vh.Type().NumIn()
		if numParam >= 1 {
			in = append(in, reflect.ValueOf(client))
		}
		if numParam == 2 {
			in = append(in, reflect.ValueOf(value))
		}
		out := vh.Call(in)
		if len(out) == 1 {
			errRaw := out[0]
			if err, safe := errRaw.Interface().(error); safe {
				return err
			}
		}
		return fmt.Errorf("when call event expected have result")
	}
	return fmt.Errorf("handler is not an function")
}
