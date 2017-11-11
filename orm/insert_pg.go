package orm

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/inconshreveable/log15"
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

func (o *ORM) postgreInsertHelper(ctx context.Context, stmt *sql.Stmt, bindArgs []interface{}, obj * object.Object, callerSuppliesPK bool, tracing bool, objTable *schema.Table) (int64, error) {
	errorString := "Insert error"
	var err error
	var lastID int64
	// TODO: CallerSuppliesPrimaryPK implementation
	if callerSuppliesPK {
		_ = stmt.QueryRowContext(ctx, bindArgs...)
	} else {
		fmt.Println("BINDARGS->",bindArgs)
		err = stmt.QueryRowContext(ctx, bindArgs...).Scan(&lastID)
	}
	if err != nil {
		if tracing {
			log15.Error(errorString, "Scan error", err)
		}
		return 0, err
	}

	// If the user supplies the primary key for this table, there is no need
	// for us to bother with populating the result of LastInsertID().
	if !callerSuppliesPK {
		if lastID != 0 {
			if tracing {
				fmt.Println("DEBUG Insert received lastID=", lastID)
			}
			obj.SetCore(objTable.Primary, lastID) // Set the new primary key in the object
		} else {
			// TODO: error condition here
			// TODO: investigate cleanup and possibility of errors in
			// current Insert function in this case
		}
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
	return 1, nil

}

