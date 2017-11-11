package orm

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/rbastic/dyndao/object"
	"github.com/inconshreveable/log15"
)

// Delete function will DELETE a record ...
func (o *ORM) Delete(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	sg := o.sqlGen
	tracing := sg.Tracing
	errorString := "Delete error"

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	objTable := o.s.GetTable(obj.Type)
	if objTable == nil {
		return 0, errors.New("Delete: unknown object table " + obj.Type)
	}

	err := o.CallBeforeDeleteHookIfNeeded(obj)
	if err != nil {
		if tracing {
			log15.Error(errorString, "BeforeUpdateHookError", err)
		}
		return 0, err
	}

	sqlStr, bindWhere, err := sg.BindingDelete(o.sqlGen, o.s, obj)
	if err != nil {
		return 0, err
	}

	if sg.Tracing {
		fmt.Printf("Delete: sqlStr->%s, bindWhere->%v\n", sqlStr, bindWhere)
	}

	stmt, err := stmtFromDbOrTx(ctx, o, tx, sqlStr)
	if err != nil {
		return 0, err
	}

	defer func() {
		stmtErr := stmt.Close()
		if stmtErr != nil {
			fmt.Println(stmtErr) // TODO: logger implementation
		}
	}()

	res, err := stmt.ExecContext(ctx, bindWhere...)
	if err != nil {
		return 0, errors.Wrap(err, "Delete")
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAff == 0 {
		return 0, ErrNoResult
	}

	err = o.CallAfterDeleteHookIfNeeded(obj)
	if err != nil {
		if tracing {
			log15.Error(errorString, "BeforeAfterUpdateHookError", err)
		}
		return 0, err
	}

	obj.MarkDirty(false)      // Flag that the object has been recently saved
	obj.ResetChangedColumns() // Reset the 'changed fields', if any

	return rowsAff, nil

}
