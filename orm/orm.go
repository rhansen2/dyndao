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
	if !objTable.MultiKey {
		_, ok := obj.KV[objTable.Primary]
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
	if !objTable.MultiKey {
		// TODO: It would be nice to capture the ID back and store it.

		// NOTE: perhaps the generator should become a part of the schema...
		// This should work well once we understand OOP in Go a bit better.
		// We should set the generators prior to running any ORM operations.
		gen := sqlitegen.New("sqlite", "test", sch)

		sql, bindArgs, err := gen.BindingInsert(obj.Type, obj.KV)
		if err != nil {
			return 0, err
		}

		//		fmt.Println(sql)
		//		fmt.Println(bindArgs)
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
		// Set our new primary key in the objcet
		obj.Set(objTable.Primary, newID)
		// Note that the object has been saved recently
		obj.SetSaved(true)
		// Reset the 'changed fields' if any exist
		obj.ResetChangedFields()

		return rowsAff, nil
	}

	panic("Insert() does not yet support MultiKey")
}

// Update function will UPDATE a record depending on various values
func Update(db *sql.DB, sch *schema.Schema, obj *object.Object) (int64, error) {
	objTable := sch.Tables[obj.Type]
	if objTable == nil {
		return 0, errors.New("Update: unknown object table " + obj.Type)
	}

	panic("Update not yet implemented")
}
