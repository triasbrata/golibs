package transaction

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triasbrata/golibs/pkg/dbx"
	mock_db "github.com/triasbrata/golibs/pkg/dbx/dbxm"
)

func TestTx(t *testing.T) {
	type args struct {
		ctx context.Context
		db  dbx.DB
		h   func(dbx.Tx) error
	}
	mockDB := mock_db.NewDB(t)
	mockTx := mock_db.NewTx(t)
	tests := []struct {
		name    string
		args    args
		wantErr error
		mock    func()
	}{
		{
			name: "success",
			mock: func() {
				mockDB.On("BeginTxx", context.Background(), &sql.TxOptions{}).Return(mockTx, nil).Once()
				mockTx.On("Commit").Return(nil).Once()
				mockTx.On("ExecContext", context.Background(), "test", 1).Return(driver.ResultNoRows, nil).Once()
			},
			args: args{
				ctx: context.Background(),
				db:  mockDB,
				h: func(tx dbx.Tx) error {
					_, err := tx.ExecContext(context.Background(), "test", 1)
					if err != nil {
						return err
					}
					return nil

				},
			},
		},
		{
			name: "fail create tx",
			mock: func() {
				mockDB.On("BeginTxx", context.Background(), &sql.TxOptions{}).Return(mockTx, fmt.Errorf("boom")).Once()
			},
			wantErr: fmt.Errorf("boom"),
			args: args{
				ctx: context.Background(),
				db:  mockDB,
				h: func(tx dbx.Tx) error {
					_, err := tx.ExecContext(context.Background(), "test", 1)
					if err != nil {
						return err
					}
					return nil

				},
			},
		},
		{
			name: "fail",
			mock: func() {
				mockDB.On("BeginTxx", context.Background(), &sql.TxOptions{}).Return(mockTx, nil).Once()
				mockTx.On("Commit").Return(nil).Once()
				mockTx.On("Rollback").Return(nil).Once()
				mockTx.On("ExecContext", context.Background(), "test", 1).Return(driver.ResultNoRows, fmt.Errorf("boom")).Once()
			},
			wantErr: fmt.Errorf("boom"),
			args: args{
				ctx: context.Background(),
				db:  mockDB,
				h: func(tx dbx.Tx) error {
					_, err := tx.ExecContext(context.Background(), "test", 1)
					if err != nil {
						return err
					}
					return nil

				},
			},
		},
		{
			name: "fail rollback",
			mock: func() {
				mockDB.On("BeginTxx", context.Background(), &sql.TxOptions{}).Return(mockTx, nil).Once()
				mockTx.On("Commit").Return(nil).Once()
				mockTx.On("Rollback").Return(fmt.Errorf("rollback err")).Once()
				mockTx.On("ExecContext", context.Background(), "test", 1).Return(driver.ResultNoRows, fmt.Errorf("boom")).Once()
			},
			wantErr: errors.Join(fmt.Errorf("boom"), fmt.Errorf("rollback err")),
			args: args{
				ctx: context.Background(),
				db:  mockDB,
				h: func(tx dbx.Tx) error {
					_, err := tx.ExecContext(context.Background(), "test", 1)
					if err != nil {
						return err
					}
					return nil

				},
			},
		},
		{
			name: "fail commit",
			mock: func() {
				mockDB.On("BeginTxx", context.Background(), &sql.TxOptions{}).Return(mockTx, nil).Once()
				mockTx.On("Commit").Return(fmt.Errorf("commit")).Once()
				mockTx.On("Rollback").Return(fmt.Errorf("rollback err")).Once()
				mockTx.On("ExecContext", context.Background(), "test", 1).Return(driver.ResultNoRows, fmt.Errorf("boom")).Once()
			},
			wantErr: errors.Join(fmt.Errorf("boom"), fmt.Errorf("rollback err"), fmt.Errorf("commit")),
			args: args{
				ctx: context.Background(),
				db:  mockDB,
				h: func(tx dbx.Tx) error {
					_, err := tx.ExecContext(context.Background(), "test", 1)
					if err != nil {
						return err
					}
					return nil

				},
			},
		},
		{
			name: "fail commit only",
			mock: func() {
				mockDB.On("BeginTxx", context.Background(), &sql.TxOptions{}).Return(mockTx, nil).Once()
				mockTx.On("Commit").Return(fmt.Errorf("commit")).Once()
				mockTx.On("Rollback").Return(nil).Once()
				mockTx.On("ExecContext", context.Background(), "test", 1).Return(driver.ResultNoRows, fmt.Errorf("boom")).Once()
			},
			wantErr: errors.Join(fmt.Errorf("boom"), fmt.Errorf("commit")),
			args: args{
				ctx: context.Background(),
				db:  mockDB,
				h: func(tx dbx.Tx) error {
					_, err := tx.ExecContext(context.Background(), "test", 1)
					if err != nil {
						return err
					}
					return nil

				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := Tx(tt.args.ctx, tt.args.db, tt.args.h)
			if tt.wantErr != nil {
				fmt.Printf("%v", err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
