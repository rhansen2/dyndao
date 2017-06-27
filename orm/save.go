package orm

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/dyndao/sqlgen/sqlitegen"
)

// TODO: For foreign key filling, we do not check to see if there are conflicts
// with regards to the uniqueness of primary key names.

func recurseAndSave(ctx context.Context, db *sql.DB, tx *sql.Tx, sch *schema.Schema, obj *object.Object) (int64, error) {
	// TODO: Implement transactions. Implement 'foreign key fill'
	// in children objects
	rowsAff, err := SaveObject(ctx, db, tx, sch, obj)
	if err != nil {
		return 0, err
	}

	table := sch.Tables[obj.Type]
	pkVal := obj.Get(table.Primary)

	// TODO: ChildrenOrder going to happen here?
	for _, v := range obj.Children {
		for _, childObj := range v {
			// set the primary key in the child object, if it exists in the child object's table
			childTable, ok := sch.Tables[childObj.Type]
			if !ok {
				return 0, errors.New("recurseAndSave: Unknown child object type " + childObj.Type + " for parent type " + obj.Type)
			}
			// check if the child schema table contains
			// the parent's primary key field as a name

			_, ok = childTable.Fields[table.Primary]
			if ok {
				// ensure it's set if so
				childObj.Set(table.Primary, pkVal)
			}

			rowsAff, err := recurseAndSave(ctx, db, tx, sch, childObj)
			if err != nil {
				return rowsAff, err
			}
		}
	}
	return rowsAff, err
}

// Save will attempt to save an entire nested object structure inside of a single transaction.
func Save(ctx context.Context, db *sql.DB, sch *schema.Schema, obj *object.Object) (int64, error) {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	// TODO: Review this code for how it uses transactions / rollbacks.
	rowsAff, err := recurseAndSave(ctx, db, tx, sch, obj)
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
		tx.Rollback()
		return 0, err
	}
	return rowsAff, nil
}

// SaveObject function will INSERT or UPDATE a record depending on
// various values.
func SaveObject(ctx context.Context, db *sql.DB, tx *sql.Tx, sch *schema.Schema, obj *object.Object) (int64, error) {
	objTable := sch.Tables[obj.Type]
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
		return Insert(ctx, db, tx, sch, obj)
	}
	return Update(ctx, db, tx, sch, obj)
}

// TODO: Read this post for more info on the above...
// https://stackoverflow.com/questions/23507531/is-golangs-sql-package-incapable-of-ad-hoc-exploratory-queries/23507765#23507765

// Insert function will INSERT a record depending on various values
func Insert(ctx context.Context, db *sql.DB, tx *sql.Tx, sch *schema.Schema, obj *object.Object) (int64, error) {
	objTable := sch.Tables[obj.Type]
	if objTable == nil {
		return 0, errors.New("Insert: unknown object table " + obj.Type)
	}
	// NOTE: perhaps the generator should become a part of the schema...
	// This should work well once we understand OOP in Go a bit better.
	// We should set the generators prior to running any ORM operations.
	gen := sqlitegen.New("sqlite", "test", sch)

	sqlStr, bindArgs, err := gen.BindingInsert(obj.Type, obj.KV)
	if err != nil {
		return 0, err
	}
	var stmt *sql.Stmt
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, sqlStr)
	} else {
		stmt, err = db.PrepareContext(ctx, sqlStr)
	}
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, bindArgs...)
	if err != nil {
		return 0, err
	}
	newID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	rowsAff, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	obj.Set(objTable.Primary, newID) // Set the new primary key in the object
	obj.SetSaved(true)               // Note that the object has been recently saved
	obj.ResetChangedFields()         // Reset the 'changed fields', if any
	return rowsAff, nil
}

// Update function will UPDATE a record depending on various values
func Update(ctx context.Context, db *sql.DB, tx *sql.Tx, sch *schema.Schema, obj *object.Object) (int64, error) {
	objTable := sch.Tables[obj.Type]
	if objTable == nil {
		return 0, errors.New("Update: unknown object table " + obj.Type)
	}
	// NOTE: perhaps the generator should become a part of the schema...
	// This should work well once I grok OOP in Go a bit better w.r.t. how this should
	// all be structured.
	// Perhaps I should set the generators prior to running any ORM operations.
	gen := sqlitegen.New("sqlite", "test", sch)

	sqlStr, bindArgs, bindWhere, err := gen.BindingUpdate(sch, obj)
	if err != nil {
		return 0, err
	}

	var stmt *sql.Stmt
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, sqlStr)
	} else {
		stmt, err = db.PrepareContext(ctx, sqlStr)
	}
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
