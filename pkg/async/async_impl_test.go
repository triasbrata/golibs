package async

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_catch(t *testing.T) {
	errChan := make(chan error)
	go func() {
		defer catch("test", errChan)
		var a map[string]interface{}
		a["test"] = "boom"
	}()
	err := <-errChan
	assert.Equal(t, fmt.Errorf("got panic when execute %s: %v", "test", "assignment to entry in nil map"), err)
}

func Test_DoWithMaxConcurrency(t *testing.T) {
	async := New()
	now := time.Now()
	arr := []int64{1, 2, 3, 4, 5}
	calledAt := []int64{}
	for _, v := range arr {
		async.Add(fmt.Sprintf("test_%v", v), func(ctx context.Context) (interface{}, error) {
			time.Sleep(100 * time.Millisecond)
			calledAt = append(calledAt, time.Since(now).Milliseconds())
			return v, nil
		})
	}
	res, err := async.DoWithMaxConcurrency(context.Background(), 2)
	assert.Len(t, calledAt, 5)
	assert.Equal(t, calledAt[0], calledAt[1])
	assert.Equal(t, calledAt[2], calledAt[3])
	assert.GreaterOrEqual(t, calledAt[4], calledAt[3])
	assert.GreaterOrEqual(t, calledAt[2], calledAt[1])
	assert.Nil(t, err)
	for _, v := range arr {
		temp, ok := res[fmt.Sprintf("test_%v", v)]
		assert.EqualValues(t, v, temp)
		assert.Equal(t, true, ok)
	}
}

func Test_Do(t *testing.T) {
	async := New()
	arr := []int64{1, 2, 3, 4, 5}
	calledAt := int64(0)
	for _, v := range arr {
		async.Add(fmt.Sprintf("test_%v", v), func(kv int64) FuncAsync {
			return func(ctx context.Context) (interface{}, error) {
				time.Sleep(100 * time.Millisecond)
				return kv, nil
			}
		}(v))
	}
	res, err := async.Do(context.Background())
	assert.GreaterOrEqual(t, int64(500), atomic.LoadInt64(&calledAt))
	assert.Nil(t, err)
	for _, v := range arr {
		temp, ok := res[fmt.Sprintf("test_%v", v)]
		assert.EqualValues(t, v, temp)
		assert.Equal(t, true, ok)
	}
}

func Test_DoWithMaxConcurrency_withError(t *testing.T) {
	async := New()
	arr := []int64{1, 2, 3, 4, 5}
	for _, v := range arr {
		async.Add(fmt.Sprintf("test_%v", v), func(kv int64) FuncAsync {
			return func(ctx context.Context) (interface{}, error) {
				time.Sleep(100 * time.Millisecond)

				if kv%2 == 0 {
					return nil, fmt.Errorf("boom %v", kv)
				}
				return kv, nil
			}
		}(v))
	}
	_, err := async.DoWithMaxConcurrency(context.Background(), 2)
	assert.NotNil(t, err)
}

func Test_Do_withAllError(t *testing.T) {
	async := New()

	for v := range 1_000 {
		async.Add(fmt.Sprintf("test_%v", v), func(kv int) FuncAsync {
			return func(ctx context.Context) (interface{}, error) {
				// time.Sleep(100 * time.Millisecond)
				return nil, fmt.Errorf("boom %v", kv)
			}
		}(v))
	}
	_, err := async.DoWithMaxConcurrency(context.Background(), 2)
	assert.NotNil(t, err)
}
