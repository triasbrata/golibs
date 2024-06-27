package dbx

import (
	"database/sql"
	"testing"

	"context"
	"fmt"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/triasbrata/golibs/go/utils"
)

func Test_dbx_Close(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	dbx := New(sqlxDB)
	expect := fmt.Errorf("boom")
	mock.ExpectClose().WillReturnError(expect)
	got := dbx.Close()
	assert.Equal(t, expect, got)
}

func Test_dbx_BeginTxx(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	dbx := New(sqlxDB)
	expect := fmt.Errorf("boom")
	mock.ExpectBegin().WillReturnError(expect)
	_, got := dbx.BeginTxx(context.Background(), &sql.TxOptions{})
	assert.Equal(t, expect, got)
}

func Test_dbx_BindNamed(t *testing.T) {
	db, _, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	dbx := New(sqlxDB)
	str, arg, err := dbx.BindNamed("select :name", utils.H{"name": 1})
	assert.Equal(t, "select ?", str)
	assert.Equal(t, []interface{}{1}, arg)
	assert.Nil(t, err)
}
func Test_dbx_ExecContext(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	dbx := New(sqlxDB)
	mock.ExpectExec("insert into ?").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
	_, err := dbx.ExecContext(context.Background(), "insert into ?", 1)
	assert.Nil(t, err)
}
func Test_dbx_GetContext(t *testing.T) {
	type Result struct {
		Id int64 `db:"id"`
	}
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	dbx := New(sqlxDB)
	query := "select * from bas where id = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	dest := []Result{}
	err := dbx.SelectContext(context.Background(), &dest, query, 1)
	assert.Nil(t, err)
	assert.Equal(t, []Result{{Id: 1}}, dest)
}
func Test_dbx_NamedExecContext(t *testing.T) {
	type Result struct {
		Id int64 `db:"id"`
	}
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	dbx := New(sqlxDB)
	query := "select * from bas where id = :id"
	arg := utils.H{"id": 1}
	q, _, _ := sqlx.BindNamed(sqlx.QUESTION, query, arg)
	mock.ExpectExec(regexp.QuoteMeta(q)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := dbx.NamedExecContext(context.Background(), query, arg)
	assert.Nil(t, err)
	lastInsert, _ := res.LastInsertId()
	assert.Equal(t, int64(1), lastInsert)
}
