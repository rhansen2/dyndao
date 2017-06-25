package orm

import (
	"database/sql"
	"errors"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/dyndao/sqlgen/sqlitegen"
)

// Save function will INSERT or UPDATE a record depending on
// various values.
func Save(db *sql.DB, sch *schema.Schema, obj *object.Object) (int64, error) {
	objTable := sch.Tables[obj.Type]
	if objTable == nil {
		return 0, errors.New("Save: unknown object table " + obj.Type)
	}
	if obj.GetSaved() {
		return 0, nil
	}
	if !objTable.MultiKey {
		_, ok := obj.KV[objTable.Fields[objTable.Primary].Name]
		if !ok {
			return Insert(db, sch, obj)
		}
		return Update(db, sch, obj)
	}

	panic("Save() does not yet support MultiKey")

	//	return 0, nil
}

// Insert function will INSERT a record depending on various values
func Insert(db *sql.DB, sch *schema.Schema, obj *object.Object) (int64, error) {
	objTable := sch.Tables[obj.Type]
	if objTable == nil {
		return 0, errors.New("Insert: unknown object table " + obj.Type)
	}
	// NOTE: perhaps the generator should become a part of the schema...
	// This should work well once we understand OOP in Go a bit better.
	// We should set the generators prior to running any ORM operations.
	gen := sqlitegen.New("sqlite", "test", sch)

	sql, bindArgs, err := gen.BindingInsert(obj.Type, obj.KV)
	if err != nil {
		return 0, err
	}

	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(bindArgs...)
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
func Update(db *sql.DB, sch *schema.Schema, obj *object.Object) (int64, error) {
	objTable := sch.Tables[obj.Type]
	if objTable == nil {
		return 0, errors.New("Update: unknown object table " + obj.Type)
	}
	// NOTE: perhaps the generator should become a part of the schema...
	// This should work well once I grok OOP in Go a bit better w.r.t. how this should
	// all be structured.
	// Perhaps I should set the generators prior to running any ORM operations.
	gen := sqlitegen.New("sqlite", "test", sch)

	sql, bindArgs, err := gen.BindingUpdate(sch, obj)
	if err != nil {
		return 0, err
	}

	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(bindArgs...)
	if err != nil {
		return 0, err
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	obj.SetSaved(true)       // Note that the object has been recently saved
	obj.ResetChangedFields() // Reset the 'changed fields', if any

	return rowsAff, nil

}
