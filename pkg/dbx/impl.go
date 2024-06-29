package dbx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type t_db struct {
	db *sqlx.DB
}

// Close implements wrappers.DB.
func (d *t_db) Close() error {
	return d.db.Close()
}

// BeginTxx implements wrappers.DB.
func (d *t_db) BeginTxx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	return d.db.BeginTxx(ctx, opts)
}

// BindNamed implements wrappers.DB.
func (d *t_db) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	return d.db.BindNamed(query, arg)
}

// ExecContext implements wrappers.DB.
func (d *t_db) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

// GetContext implements wrappers.DB.
func (d *t_db) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.db.GetContext(ctx, dest, query, args...)
}

// NamedExecContext implements wrappers.DB.
func (d *t_db) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return d.db.NamedExecContext(ctx, query, arg)
}

// NamedQuery implements wrappers.DB.
func (d *t_db) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return d.db.NamedQuery(query, arg)
}

// NamedQueryContext implements wrappers.DB.
func (d *t_db) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	return d.db.NamedQueryContext(ctx, query, arg)
}

// SelectContext implements wrappers.DB.
func (d *t_db) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.db.SelectContext(ctx, dest, query, args...)
}

func New(db *sqlx.DB) DB {
	return &t_db{db: db}
}
