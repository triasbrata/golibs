package async

import (
	"context"

	"github.com/triasbrata/golibs/pkg/utils"
)

type FuncAsync = func(ctx context.Context) (interface{}, error)

type Async interface {
	Add(funcName string, f FuncAsync)
	Do(ctx context.Context) (utils.H, error)
	DoWithMaxConcurrency(ctx context.Context, maxConcurrency int) (utils.H, error)
}
