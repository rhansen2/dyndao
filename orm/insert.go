package orm

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// Insert function will INSERT a record, given an optional transaction and an object.
// It returns the number of rows affected (int64) and any error that may have occurred.
func (o *ORM) Insert(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	sg := o.sqlGen
	tracing := sg.Tracing
	errorString := "Insert error"

	// Check context
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	// Lookup schema table
	objTable := o.s.GetTable(obj.Type)
	if objTable == nil {
		if tracing {
			log15.Error(errorString, "GetTable_error", "objTable was unknown")
		}
		return 0, errors.New("Insert: unknown object table " + obj.Type)
	}

	callerSuppliesPK := objTable.CallerSuppliesPK

	// Call any before create hooks
	err := o.CallBeforeCreateHookIfNeeded(obj)
	if err != nil {
		if tracing {
			log15.Error(errorString, "BeforeCreateHookError", err)
		}
		return 0, err
	}

	// Prepare our binding insert SQL statement and the binding parameters
	sqlStr, bindArgs, err := sg.BindingInsert(sg, o.s, obj.Type, obj.KV)
	if err != nil {
		if tracing {
			log15.Error(errorString, "BindingInsert_error", err)
		}
		return 0, err
	}
	if tracing {
		fmt.Println("Insert/sqlStr=", sqlStr, "bindArgs=", bindArgs)
	}

	// Potential way to capture LastInsertID
	// Oracle-specific fix
	var lastID int64
	if (!callerSuppliesPK) && o.sqlGen.FixLastInsertIDbug && sg.IsORACLE {
		bindArgs = append(bindArgs, sql.Named(o.s.GetTable(obj.Type).Primary, sql.Out{
			Dest: &lastID,
		}))
	}

	// Prepare statement handle from either the database or the transaction
	stmt, err := stmtFromDbOrTx(ctx, o, tx, sqlStr)
	if err != nil {
		if tracing {
			log15.Error(errorString, "stmtFromDbOrTx_error", err)
		}

		return 0, err
	}
	// TODO: Check if we can replace maybeDereferenceArgs now
	newBindArgs := make([]interface{}, len(bindArgs))
	for i, arg := range bindArgs {
		newBindArgs[i] = maybeDereferenceArgs(arg)
	}

	if sg.IsPOSTGRES || sg.IsDB2 {
		return o.postgreInsertHelper(ctx, stmt, bindArgs, obj, callerSuppliesPK, tracing, objTable)
	}
	return o.insertHelper(ctx, stmt, bindArgs, obj, callerSuppliesPK, tracing, objTable, &lastID)
}

func (o *ORM) insertHelper(ctx context.Context, stmt *sql.Stmt, bindArgs []interface{}, obj * object.Object, callerSuppliesPK bool, tracing bool, objTable *schema.Table, lastID *int64) (int64, error) {
	errorString := "Insert error"
	var err error
	defer func() {
		//fmt.Println("DEFER INSERT ABOUT TO CLOSE")
		err := stmt.Close()
		if err != nil {
			fmt.Println("DEFER INSERT ERROR stmt.Close error=", err) // TODO: logging implementation
			return
		}
		//fmt.Println("DEFER INSERT CLOSED")
	}()

	// Execute our statement
	res, err := stmt.ExecContext(ctx, bindArgs...)
	if err != nil {
		if tracing {
			log15.Error(errorString, "ExecContext_error", err)
			fmt.Println("orm/save error", err)
		}

		return 0, errors.Wrap(err, "Insert/ExecContext")
	}

	// If the user supplies the primary key for this table, there is no need
	// for us to bother with populating the result of LastInsertID().
	if !callerSuppliesPK {
		newID, err := res.LastInsertId()
		if err != nil && *lastID == 0 {
			if tracing {
				fmt.Println("orm/save error", err)
			}
			log15.Error(errorString, "LastInsertID_error", err)
			return 0, err
		}
		if *lastID != 0 {
			newID = *lastID
		}
		if tracing {
			fmt.Println("DEBUG Insert received newID=", newID)
		}

		obj.SetCore(objTable.Primary, newID) // Set the new primary key in the object
	}

	// Check rows affected
	rowsAff, err := res.RowsAffected()
	if err != nil {
		if tracing {
			fmt.Println("orm/save error", err)
		}
		return 0, err
	}

	if rowsAff == 0 {
		return 0, ErrNoResult
	}

	// Call after create hook
	err = o.CallAfterCreateHookIfNeeded(obj)
	if err != nil {
		if tracing {
			log15.Error(errorString, "BeforeAfterCreateHookError", err)
		}
		return 0, err
	}

	obj.MarkDirty(false)      // Note that the object has been recently saved
	obj.ResetChangedColumns() // Reset the 'changed fields', if any
	return rowsAff, nil
}

