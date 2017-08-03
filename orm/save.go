package orm

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// NOTE: For foreign key filling, we do not check to see if there are conflicts
// with regards to the uniqueness of primary key names.
func pkQueryValsFromKV(obj *object.Object, sch *schema.Schema, parentTableName string) (map[string]interface{}, error) {
	qv := make(map[string]interface{})

	schemaTable := sch.GetTable(parentTableName)

	if schemaTable == nil {
		return nil, errors.New("pkQueryValsFromKV: no schemaTable for table " + parentTableName)
	}
	schemaPrimary := schemaTable.Primary

	for fName, field := range schemaTable.Fields {
		if field.IsIdentity || field.IsForeignKey || field.Name == schemaPrimary {
			qv[fName] = obj.Get(fName)
		}
	}
	return qv, nil
}

// TODO: add a queryVals struct containing all new LastInsertIDs
func (o ORM) recurseAndSave(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	rowsAff, err := o.SaveObject(ctx, tx, obj)
	if err != nil {
		return 0, err
	}

	table := o.s.GetTable(obj.Type)
	pkVal := obj.Get(table.Primary)

	for _, v := range obj.Children {
		for _, childObj := range v {
			// set the primary key in the child object, if it exists in the child object's table
			childTable, ok := o.s.Tables[childObj.Type]
			if !ok {
				return 0, errors.New("recurseAndSave: Unknown child object type " + childObj.Type + " for parent type " + obj.Type)
			}
			// TODO: support propagation of additional primary keys that are
			// saved from previous recursive saves
			// ... check if the child schema table contains
			// the parent's primary key field as a name
			_, ok = childTable.Fields[table.Primary]
			if ok {
				// set in the child object if the table contains the primary
				childObj.Set(table.Primary, pkVal)
			}

			aff, err := o.recurseAndSave(ctx, tx, childObj)
			if err != nil {
				return rowsAff + aff, err
			}
		}
	}
	return rowsAff, err
}

// SaveAll will attempt to save an entire nested object structure inside of a single transaction.
// It begins the transaction, attempts to recursively save the object and all of it's children,
// and any of the children's children, and then will finally rollback/commit as necessary.
func (o ORM) SaveAll(ctx context.Context, obj *object.Object) (int64, error) {
	tx, err := o.RawConn.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	// TODO: Review this code for how it uses transactions / rollbacks.
	rowsAff, err := o.recurseAndSave(ctx, tx, obj)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			// TODO: Not sure if this wrap is right.
			return 0, errors.Wrap(err, err2.Error())
		}
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		// TODO: Error is not checked.
		tx.Rollback()
		return 0, err
	}
	return rowsAff, nil
}

// SaveObject function will INSERT or UPDATE a record. It does not attempt to
// save any of the children. If given a transaction, it will use that to
// attempt to insert the data.
func (o ORM) SaveObject(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	objTable := o.s.GetTable(obj.Type)
	if objTable == nil {
		return 0, errors.New("SaveObject: unknown object table " + obj.Type)
	}
	if obj.GetSaved() {
		return 0, nil
	}
	fieldMap := objTable.Fields
	pk := objTable.Primary
	if pk == "" {
		return 0, errors.New("SaveObject: empty primary key for " + obj.Type)
	}
	f := fieldMap[pk]
	if f == nil {
		return 0, errors.New("SaveObject: empty field " + pk + " for " + obj.Type)
	}
	// Check the primary key to see if we should insert or update
	_, ok := obj.KV[f.Name]
	if !ok {
		return o.Insert(ctx, tx, obj)
	}
	return o.Update(ctx, tx, obj)
}

func stmtFromDbOrTx(ctx context.Context, o ORM, tx *sql.Tx, sqlStr string) (*sql.Stmt, error) {
	var stmt *sql.Stmt
	var err error
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, sqlStr)
	} else {
		stmt, err = o.RawConn.PrepareContext(ctx, sqlStr)
	}
	return stmt, err
}

// Insert function will INSERT a record depending on various values
func (o ORM) Insert(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	objTable := o.s.GetTable(obj.Type)
	if objTable == nil {
		if os.Getenv("DEBUG") != "" {
			log15.Info("orm/save error", "error", "thing was unknown")
		}
		return 0, errors.New("Insert: unknown object table " + obj.Type)
	}
	sqlStr, bindArgs, err := o.sqlGen.BindingInsert(o.s, obj.Type, obj.KV)
	if err != nil {
		if os.Getenv("DEBUG") != "" {
			log15.Info("orm/save error", "error", err)
		}
		return 0, err
	}
	if os.Getenv("DEBUG") != "" {
		fmt.Println("Insert/sqlStr=", sqlStr, "bindArgs=", bindArgs)
	}

	// FIXME: Possible bug in rana ora.v4? I wouldn't have expected that I'd
	// have to append a parameter like this, based on reading the code.
	if !o.sqlGen.CallerSuppliesPrimaryKey() {
		if o.sqlGen.FixLastInsertIDbug() {
			var lastID int64
			bindArgs = append(bindArgs, &lastID)
		}
	} /*else {
				var lastID string
				bindArgs = append(bindArgs, &lastID)
	}*/
	stmt, err := stmtFromDbOrTx(ctx, o, tx, sqlStr)
	if err != nil {
		if os.Getenv("DEBUG") != "" {
			log15.Info("orm/save error", "error", err)
		}

		return 0, err
	}
	// TODO: Error return value not checked.
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, bindArgs...)
	if err != nil {
		if os.Getenv("DEBUG") != "" {
			fmt.Println("orm/save error", err)
		}

		return 0, errors.Wrap(err, "Insert/ExecContext")
	}

	// If we are not expecting the caller to supply the primary key,
	// then we should not try to capture the last value (for example,
	// using LAST_INSERT_ID() with MySQL..)
	// TODO: Should CallerSuppliesPrimaryKey be per-table?
	if !o.sqlGen.CallerSuppliesPrimaryKey() {
		newID, err := res.LastInsertId()
		if err != nil {

			if os.Getenv("DEBUG") != "" {
				fmt.Println("orm/save error", err)
			}
			return 0, err
		}
		if os.Getenv("DEBUG") != "" {
			fmt.Println("DEBUG Insert received newID=", newID)
		}

		obj.Set(objTable.Primary, newID) // Set the new primary key in the object
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {

		if os.Getenv("DEBUG") != "" {
			fmt.Println("orm/save error", err)
		}
		return 0, err
	}
	obj.SetSaved(true)       // Note that the object has been recently saved
	obj.ResetChangedFields() // Reset the 'changed fields', if any
	return rowsAff, nil
}

// Update function will UPDATE a record ...
func (o ORM) Update(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	objTable := o.s.GetTable(obj.Type)

	if objTable == nil {
		return 0, errors.New("Update: unknown object table " + obj.Type)
	}
	sqlStr, bindArgs, bindWhere, err := o.sqlGen.BindingUpdate(o.s, obj)
	if err != nil {
		return 0, err
	}
	if os.Getenv("DEBUG") != "" {
		fmt.Println("Update/sqlStr=", sqlStr, "bindArgs=", bindArgs, "bindWhere=", bindWhere)
	}

	stmt, err := stmtFromDbOrTx(ctx, o, tx, sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	allBind := append(bindArgs, bindWhere...)
	res, err := stmt.ExecContext(ctx, allBind...)
	if err != nil {
		return 0, errors.Wrap(err, "Update")
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	obj.SetSaved(true)       // Note that the object has been recently saved
	obj.ResetChangedFields() // Reset the 'changed fields', if any

	return rowsAff, nil
}
