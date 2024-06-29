package async

import (
	"context"
	"fmt"
	"sync"

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
	wg := sync.WaitGroup{}
	lenFunc := len(t.funcHolder)
	wg.Add(lenFunc)
	var sem chan struct{}
	if maxConcurrency > 0 {
		sem = make(chan struct{}, maxConcurrency)
	}
	errChan := make(chan error, lenFunc)
	for key, fu := range t.funcHolder {
		if maxConcurrency > 0 {
			sem <- struct{}{}
		}
		if call, safe := fu.(FuncAsync); safe {
			go func() {
				defer func() {
					if maxConcurrency > 0 {
						<-sem
					}
				}()
				defer wg.Done()
				defer catch(key, errChan)

				resFunc, er := call(ctx)
				if er != nil {
					fmt.Printf("err: %v\n", er)
					errChan <- er
					return
				}
				res[key] = resFunc
			}()
		}
	}
	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return make(map[string]interface{}), err
	}
	close(errChan)
	return res, err
}

func New() Async {
	return &ta{
		funcHolder: make(utils.H),
	}
}

func catch(funcCaller string, err chan error) {
	if r := recover(); r != nil {
		err <- fmt.Errorf("got panic when execute %s: %v", funcCaller, r)
	}
}
