package orm

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/rbastic/dyndao/object"
)

// Delete function will DELETE a record ...
func (o ORM) Delete(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	objTable := o.s.GetTable(obj.Type)
	if objTable == nil {
		return 0, errors.New("Delete: unknown object table " + obj.Type)
	}
	sqlStr, bindWhere, err := o.sqlGen.BindingDelete(o.sqlGen, o.s, obj)
	if err != nil {
		return 0, err
	}
	if os.Getenv("DEBUG") != "" {
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

	obj.SetSaved(true)       // Flag that the object has been recently saved
	obj.ResetChangedColumns() // Reset the 'changed fields', if any

	return rowsAff, nil

}
