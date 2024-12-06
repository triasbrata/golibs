package async

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/triasbrata/golibs/pkg/utils"
	"golang.org/x/sync/errgroup"
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
	group, ctx := errgroup.WithContext(ctx)
	if maxConcurrency > 0 {
		group.SetLimit(maxConcurrency)
	}
	res = utils.H{}
	resMap := &sync.Map{}
	// if have max concurancy then we will create an concuracy controller with semantic mecanism
	for key, fu := range t.funcHolder {
		if call, safe := fu.(FuncAsync); safe {
			group.Go(func(inkey string) func() error {
				return func() (err error) {
					defer func() {
						if r := recover(); r != nil {
							errPanic := fmt.Errorf("%+v\n\t%s", r, debug.Stack())
							if err != nil {
								err = errors.Join(err, errPanic)
							} else {
								err = errPanic
							}
						}
					}()
					resFunc, err := call(ctx)
					if err != nil {
						return err
					}
					resMap.Store(inkey, resFunc)
					return nil
				}
			}(key))
		}
	}
	err = group.Wait()
	if err != nil {
		return res, err
	}
	resMap.Range(func(key, value any) bool {
		skey := key.(string)
		res[skey] = value
		return true
	})
	return res, nil
}

func New() Async {
	return &ta{
		funcHolder: make(utils.H),
	}
}
