package orm

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/rbastic/dyndao/object"
)

// Delete function will DELETE a record ...
func (o ORM) Delete(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	objTable := o.s.Tables[obj.Type]
	if objTable == nil {
		return 0, errors.New("Delete: unknown object table " + obj.Type)
	}
	sqlStr, bindArgs, bindWhere, err := o.sqlGen.BindingDelete(o.s, obj)
	if err != nil {
		return 0, err
	}

	stmt, err := stmtFromDbOrTx(ctx, o, tx, sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	fmt.Println("WARN bindArgs->", bindArgs)
	fmt.Println("WARN bindWhere->", bindWhere)

	allBind := append(bindArgs, bindWhere...)
	res, err := stmt.ExecContext(ctx, allBind...)
	if err != nil {
		return 0, errors.Wrap(err, "Delete")
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	obj.SetSaved(true)       // Note that the object has been recently saved
	obj.ResetChangedFields() // Reset the 'changed fields', if any

	return rowsAff, nil

}
