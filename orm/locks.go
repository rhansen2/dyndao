package orm

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"fmt"

)

func (o *ORM) GetLock(ctx context.Context, tx *sql.Tx, lockStr string) (bool, error) {
	sg := o.sqlGen
	//tracing := sg.Tracing
	//errorString := "GetLock error"

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	sqlStr, bindWhere, err := sg.GetLock(o.sqlGen, o.s, lockStr)
	if err != nil {
		return false, err
	}

	if sg.Tracing {
		fmt.Printf("GetLock: sqlStr->%s, bindWhere->%v\n", sqlStr, bindWhere)
	}

	stmt, err := stmtFromDbOrTx(ctx, o, tx, sqlStr)
	if err != nil {
		return false, err
	}

	defer func() {
		stmtErr := stmt.Close()
		if stmtErr != nil {
			fmt.Println(stmtErr) // TODO: logger implementation
		}
	}()

	res, err := stmt.ExecContext(ctx, bindWhere...)
	if err != nil {
		return false, errors.Wrap(err, "GetLock")
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAff == 0 {
		return false, ErrNoResult
	}

	return true, nil
}

func (o *ORM) ReleaseLock(ctx context.Context, tx *sql.Tx, lockStr string) (bool, error) {
	sg := o.sqlGen

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	sqlStr, bindWhere, err := sg.ReleaseLock(o.sqlGen, o.s, lockStr)
	if err != nil {
		return false, errors.Wrap(err, "ReleaseLock/sg")
	}

	if sg.Tracing {
		fmt.Printf("ReleaseLock: sqlStr->%s, bindWhere->%v\n", sqlStr, bindWhere)
	}

	stmt, err := stmtFromDbOrTx(ctx, o, tx, sqlStr)
	if err != nil {
		return false, errors.Wrap(err, "ReleaseLock/stmtFromDbOrTx")
	}

	defer func() {
		stmtErr := stmt.Close()
		if stmtErr != nil {
			fmt.Println(stmtErr) // TODO: logger implementation
		}
	}()

	res, err := stmt.ExecContext(ctx, bindWhere...)
	if err != nil {
		return false, errors.Wrap(err, "ReleaseLock/ExecContext")
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "ReleaseLock/RowsAffected")
	}

	if rowsAff == 0 {
		return false, ErrNoResult
	}

	return true, nil
}
