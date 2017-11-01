// Package dao is the data access object swiss army knife / "black box".
package orm

import (
	"context"
	"database/sql"

	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	//"github.com/rbastic/dyndao/
	"fmt"
	"runtime/debug"
)

type TxFuncType func(*sql.Tx) error

// Transact is meant to group operations into transactions, simplify error
// handling, and recover from any panics.  See:
// http://stackoverflow.com/questions/16184238/database-sql-tx-detecting-commit-or-rollback
// Please note this function has been changed from the above post to use
// contexts
func (o *ORM) Transact(ctx context.Context, txFunc TxFuncType, opts *sql.TxOptions) error {
	tx, err := o.RawConn.BeginTx(ctx, opts)
	if err != nil {
		log15.Error("[Transact]", "BeginTx", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}

			err = fmt.Errorf("%s [Transact/defer/panic %s]", err, debug.Stack())
		}
		if err != nil {
			rollbackErr := tx.Rollback()

			if rollbackErr != nil {
				err = errors.Wrap(err, rollbackErr.Error())
			}
			return
		}
		err2 := tx.Commit()
		if err2 != nil {
			err = errors.Wrap(err, err2.Error())
		}
	}()

	err = txFunc(tx)

	if err != nil {
		return err
	}

	return nil
}

func (o *ORM) TransactRethrow(ctx context.Context, txFunc TxFuncType, opts *sql.TxOptions) error {
	tx, err := o.RawConn.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}
