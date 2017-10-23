package orm

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/object"
)

// Update function will UPDATE a record ...
func (o ORM) Update(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	errorString := "Update error"

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	err := o.CallBeforeUpdateHookIfNeeded(obj)
	if err != nil {
		if os.Getenv("DB_TRACE") != "" {
			log15.Error(errorString, "BeforeUpdateHookError", err)
		}
		return 0, err
	}

	sg := o.sqlGen
	sqlStr, bindArgs, bindWhere, err := sg.BindingUpdate(sg, o.s, obj)
	if err != nil {
		if os.Getenv("DB_TRACE") != "" {
			fmt.Println("Update/sqlStr, err=", err)
		}
		return 0, err
	}
	if os.Getenv("DB_TRACE") != "" {
		fmt.Println("Update/sqlStr=", sqlStr, "bindArgs=", bindArgs, "bindWhere=", bindWhere)
	}

	stmt, err := stmtFromDbOrTx(ctx, o, tx, sqlStr)
	if err != nil {
		return 0, err
	}
	defer func() {
		//fmt.Println("DEFER UPDATE ABOUT TO CLOSE")
		err := stmt.Close()
		if err != nil {
			fmt.Println("DEFER UPDATE ERROR stmt.Close error=", err) // TODO: logging implementation
			return
		}
		//fmt.Println("DEFER UPDATE CLOSED")
	}()

	allBind := append(bindArgs, bindWhere...)
	newAllBind := make([]interface{}, len(allBind))
	// TODO: Still necessary?
	for i, arg := range allBind {
		newAllBind[i] = maybeDereferenceArgs(arg)
	}
	res, err := stmt.ExecContext(ctx, newAllBind...)
	if err != nil {
		return 0, errors.Wrap(err, "Update")
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	err = o.CallAfterUpdateHookIfNeeded(obj)
	if err != nil {
		if os.Getenv("DB_TRACE") != "" {
			log15.Error(errorString, "BeforeAfterUpdateHookError", err)
		}
		return 0, err
	}

	obj.SetSaved(true)       // Note that the object has been recently saved
	obj.ResetChangedFields() // Reset the 'changed fields', if any

	return rowsAff, nil
}
