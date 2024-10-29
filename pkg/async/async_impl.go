package async

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/triasbrata/golibs/pkg/utils"
)

// implemenet from async
type ta struct {
	funcHolder utils.H
}

// Add implements Async.
func (t *ta) Add(funcName string, f func(ctx context.Context) (interface{}, error)) {
	t.funcHolder[funcName] = f
}

// Do implements Async.
func (t *ta) Do(ctx context.Context) (map[string]interface{}, error) {
	return t.do(ctx, 0)
}

// DoWithMaxConcurrency implements Async.
func (t *ta) DoWithMaxConcurrency(ctx context.Context, maxConcurrency int) (map[string]interface{}, error) {
	return t.do(ctx, maxConcurrency)
}

func (t *ta) do(ctx context.Context, maxConcurrency int) (res map[string]interface{}, err error) {
	res = utils.H{}
	resMap := &sync.Map{}
	lenFunc := len(t.funcHolder)
	counter := &atomic.Int32{}
	var sem chan struct{}
	// if have max concurancy then we will create an concuracy controller with semantic mecanism
	if maxConcurrency > 0 {
		sem = make(chan struct{}, maxConcurrency)
	}
	errChan := make(chan error, lenFunc)
	defer close(errChan)
	validReq := &atomic.Bool{}
	validReq.Store(true)
	for key, fu := range t.funcHolder {
		if maxConcurrency > 0 {
			sem <- struct{}{}
		}
		if call, safe := fu.(FuncAsync); safe {
			go func(inkey string) {
				defer func() {
					if maxConcurrency > 0 {
						<-sem
					}
					counter.Add(1)
				}()
				defer catch(inkey, errChan)

				resFunc, resErr := call(ctx)
				if resErr != nil && validReq.Load() {
					validReq.Store(false)
					errChan <- resErr
					return
				}
				resMap.Store(inkey, resFunc)
			}(key)
		}
	}

	for {
		select {
		case errTw := <-errChan:
			return res, errTw
		default:
			if counter.Load() == int32(lenFunc) {
				resMap.Range(func(key, value any) bool {
					keyString, safe := key.(string)
					if safe {
						res[keyString] = value
					}
					return true
				})
				return res, nil
			}
		}
	}

}

func New() Async {
	return &ta{
		funcHolder: make(utils.H),
	}
}

func catch(funcCaller string, err chan error) {
	if r := recover(); r != nil {
		stack := string(debug.Stack())
		stackArr := strings.Split(stack, "\n")
		newStack := make([]string, 0, len(stackArr)-5)
		for i, stack := range stackArr {
			if i >= 2 && i <= 6 {
				continue
			}
			newStack = append(newStack, stack)
		}

		strStack := strings.Join(newStack, "\n")
		err <- fmt.Errorf("got panic when execute %s: %v \n %s", funcCaller, r, strStack)
	}
}
