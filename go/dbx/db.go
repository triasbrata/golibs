package dbx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type FuncTx func(Tx) error

//go:generate mockery  --name=Tx
type Tx interface {
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	Commit() error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	Rollback() error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

//go:generate mockery  --name=DB
type DB interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	Close() error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
