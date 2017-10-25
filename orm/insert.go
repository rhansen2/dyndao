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

// Insert function will INSERT a record, given an optional transaction and an object.
// It returns the number of rows affected (int64) and any error that may have occurred.
func (o ORM) Insert(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	errorString := "Insert error"

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	objTable := o.s.GetTable(obj.Type)
	if objTable == nil {
		if os.Getenv("DB_TRACE") != "" {
			log15.Error(errorString, "GetTable_error", "thing was unknown")
		}
		return 0, errors.New("Insert: unknown object table " + obj.Type)
	}

	callerSuppliesPK := objTable.CallerSuppliesPK

	err := o.CallBeforeCreateHookIfNeeded(obj)
	if err != nil {
		if os.Getenv("DB_TRACE") != "" {
			log15.Error(errorString, "BeforeCreateHookError", err)
		}
		return 0, err
	}

	sg := o.sqlGen
	sqlStr, bindArgs, err := sg.BindingInsert(sg, o.s, obj.Type, obj.KV)
	if err != nil {
		if os.Getenv("DB_TRACE") != "" {
			log15.Error(errorString, "BindingInsert_error", err)
		}
		return 0, err
	}
	if os.Getenv("DB_TRACE") != "" {
		fmt.Println("Insert/sqlStr=", sqlStr, "bindArgs=", bindArgs)
	}

	var lastID int64
	// Oracle-specific fix
	if (!callerSuppliesPK) && o.sqlGen.FixLastInsertIDbug {
		bindArgs = append(bindArgs, sql.Named(o.s.GetTable(obj.Type).Primary, sql.Out{
			Dest: &lastID,
		}))
	}

	stmt, err := stmtFromDbOrTx(ctx, o, tx, sqlStr)
	if err != nil {
		if os.Getenv("DB_TRACE") != "" {
			log15.Error(errorString, "stmtFromDbOrTx_error", err)
		}

		return 0, err
	}
	defer func() {
		//fmt.Println("DEFER INSERT ABOUT TO CLOSE")
		err := stmt.Close()
		if err != nil {
			fmt.Println("DEFER INSERT ERROR stmt.Close error=", err) // TODO: logging implementation
			return
		}
		//fmt.Println("DEFER INSERT CLOSED")
	}()

	// TODO: Still necessary?
	newBindArgs := make([]interface{}, len(bindArgs))
	for i, arg := range bindArgs {
		newBindArgs[i] = maybeDereferenceArgs(arg)
	}

	res, err := stmt.ExecContext(ctx, bindArgs...)
	if err != nil {
		if os.Getenv("DB_TRACE") != "" {
			log15.Error(errorString, "ExecContext_error", err)
			fmt.Println("orm/save error", err)
		}

		return 0, errors.Wrap(err, "Insert/ExecContext")
	}

	// If the user supplies the primary key for this table, there is no need
	// for us to bother with the LastInsertId() check.
	if !callerSuppliesPK {
		newID, err := res.LastInsertId()
		if err != nil && lastID == 0 {
			if os.Getenv("DB_TRACE") != "" {
				fmt.Println("orm/save error", err)
			}
			log15.Error(errorString, "LastInsertID_error", err)
			return 0, err
		}
		if lastID != 0 {
			newID = lastID
		}
		if os.Getenv("DB_TRACE") != "" {
			fmt.Println("DEBUG Insert received newID=", newID)
		}

		obj.SetCore(objTable.Primary, newID) // Set the new primary key in the object
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		if os.Getenv("DB_TRACE") != "" {
			fmt.Println("orm/save error", err)
		}
		return 0, err
	}

	err = o.CallAfterCreateHookIfNeeded(obj)
	if err != nil {
		if os.Getenv("DB_TRACE") != "" {
			log15.Error(errorString, "BeforeAfterCreateHookError", err)
		}
		return 0, err
	}

	obj.SetSaved(true)       // Note that the object has been recently saved
	obj.ResetChangedColumns() // Reset the 'changed fields', if any
	return rowsAff, nil
}
