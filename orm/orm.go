package orm

import (
	"context"
	"database/sql"
	"log"

	"github.com/pkg/errors"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/dyndao/sqlgen/sqlitegen"
)

// Save function will INSERT or UPDATE a record depending on
// various values.
func Save(ctx context.Context, db *sql.DB, sch *schema.Schema, obj *object.Object) (int64, error) {
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
			return Insert(ctx, db, sch, obj)
		}
		return Update(ctx, db, sch, obj)
	}

	panic("Save() does not yet support MultiKey")

	//	return 0, nil
}

// Retrieve function will fleshen an entire object structure, given some primary keys.
func Retrieve(ctx context.Context, db *sql.DB, sch *schema.Schema, table string, pkValues map[string]interface{}) (*object.Object, error) {
	objTable := sch.Tables[table]
	if objTable == nil {
		return nil, errors.New("Retrieve: unknown object table " + table)
	}
	obj := object.New(table)
	obj.KV = pkValues

	gen := sqlitegen.New("sqlite", "test", sch)
	sql, bindArgs, err := gen.BindingRetrieve(sch, obj)
	if err != nil {
		return nil, err
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.QueryContext(ctx, bindArgs...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	if err != nil {
		return nil, err
	}
	for res.Next() {
		columnNames, err := res.Columns()
		if err != nil {
			return nil, err
		}
		columnTypes, err := res.ColumnTypes()
		if err != nil {
			return nil, err
		}

		columnPointers := make([]interface{}, len(columnNames))
		for i := 0; i < len(columnNames); i++ {

			ct := columnTypes[i]
			// TODO: Improve database type support.
			if ct.DatabaseTypeName() == "text" {
				var s string
				columnPointers[i] = &s
			} else {
				var j int64
				columnPointers[i] = &j
			}
			//			columnPointers[i] = &columns[i]
		}

		if err := res.Scan(columnPointers...); err != nil {
			log.Fatalln(err)
		}

		for i, v := range columnPointers {
			ct := columnTypes[i]
			// TODO: Improve database type support.
			if ct.DatabaseTypeName() == "text" {
				val := v.(*string)
				obj.Set(columnNames[i], *val)
			} else {
				val := v.(*int64)
				obj.Set(columnNames[i], *val)
			}
		}
	}

	err = res.Err()
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// TODO: Read this post for more info on the above... https://stackoverflow.com/questions/23507531/is-golangs-sql-package-incapable-of-ad-hoc-exploratory-queries/23507765#23507765

// Insert function will INSERT a record depending on various values
func Insert(ctx context.Context, db *sql.DB, sch *schema.Schema, obj *object.Object) (int64, error) {
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

	stmt, err := db.PrepareContext(ctx, sql)
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
func Update(ctx context.Context, db *sql.DB, sch *schema.Schema, obj *object.Object) (int64, error) {
	objTable := sch.Tables[obj.Type]
	if objTable == nil {
		return 0, errors.New("Update: unknown object table " + obj.Type)
	}
	// NOTE: perhaps the generator should become a part of the schema...
	// This should work well once I grok OOP in Go a bit better w.r.t. how this should
	// all be structured.
	// Perhaps I should set the generators prior to running any ORM operations.
	gen := sqlitegen.New("sqlite", "test", sch)

	sql, bindArgs, bindWhere, err := gen.BindingUpdate(sch, obj)
	if err != nil {
		return 0, err
	}

	stmt, err := db.PrepareContext(ctx, sql)
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
