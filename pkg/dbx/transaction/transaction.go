package transaction

import (
	"context"
	"database/sql"
	"errors"

	"github.com/triasbrata/golibs/pkg/dbx"
)

func Tx(ctx context.Context, db dbx.DB, h func(dbx.Tx) error) (err error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		errCommit := tx.Commit()
		if errCommit != nil {
			if uw, ok := err.(interface{ Unwrap() []error }); ok {
				errs := uw.Unwrap()
				errs = append(errs, errCommit)
				err = errors.Join(errs...)
			} else {
				err = errors.Join(err, errCommit)
			}
		}
	}()
	err = h(tx)
	if err != nil {
		errRolback := tx.Rollback()
		if errRolback != nil {
			err = errors.Join(err, errRolback)
		}
		return err
	}
	return err

}
