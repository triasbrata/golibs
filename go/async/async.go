package golib

import (
	"context"

	"github.com/triasbrata/golibs/go/types"
)

type FuncAsync = func(ctx context.Context) (interface{}, error)
type Async interface {
	Add(funcName string, f FuncAsync)
	Do(ctx context.Context) (types.H, error)
	DoWithMaxConcurrency(ctx context.Context, maxConcurrency int) (types.H, error)
}
