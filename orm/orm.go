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

// RetrieveParentViaChild function
func RetrieveParentViaChild(ctx context.Context, db *sql.DB, sch *schema.Schema, table string, queryValues map[string]interface{}, childObj *object.Object) (*object.Object, error) {
	objTable := sch.Tables[table]
	if objTable == nil {
		return nil, errors.New("RetrieveWithChildren: unknown object table " + table)
	}

	obj, err := RetrieveObject(ctx, db, sch, table, queryValues)
	if err != nil {
		return nil, err
	}
	// TODO: support multiple objects...
	if childObj != nil {
		if obj.Children[childObj.Type] == nil {
			obj.Children[childObj.Type] = make(object.ObjectArray, 1)
		}
		obj.Children[childObj.Type][0] = childObj
	}

	// TODO: Not sure that this approach ends up very
	// practical.
	var parentObj *object.Object
	if objTable.ParentTables != nil {
		for _, parentName := range objTable.ParentTables {
			parentObj, err = RetrieveParentViaChild(ctx, db, sch, parentName, queryValues, obj)
			if err != nil {
				return nil, err
			}
			if parentObj.Children[table] == nil {
				parentObj.Children[table] = make(object.ObjectArray, 1)
			}
			parentObj.Children[table][0] = obj
		}
	}
	if parentObj == nil {
		parentObj = obj
	}

	return parentObj, nil
}

// RetrieveWithChildren function will fleshen an *entire* object structure, given some primary keys
func RetrieveWithChildren(ctx context.Context, db *sql.DB, sch *schema.Schema, table string, pkValues map[string]interface{}) (*object.Object, error) {
	objTable := sch.Tables[table]
	if objTable == nil {
		return nil, errors.New("RetrieveWithChildren: unknown object table " + table)
	}
	obj := object.New(table)

	obj, err := RetrieveObject(ctx, db, sch, table, pkValues)
	if err != nil {
		return nil, errors.Wrap(err, "RetrieveWithChildren/RetrieveObject")
	}

	for name := range objTable.Children {
		childObj := object.New(name)
		childPkValues := make(map[string]interface{})

		childSchemaTable := sch.Tables[name]

		pVal, ok := pkValues[childSchemaTable.Primary]
		if ok {
			childPkValues[childSchemaTable.Primary] = pVal
		}

		if childSchemaTable.MultiKey && childSchemaTable.ForeignKeys != nil {
			for _, fk := range childSchemaTable.ForeignKeys {
				childPkValues[fk] = pkValues[fk]
			}
		}
		// TODO: Should we do anything else with pkValues?

		childObj, err := RetrieveObject(ctx, db, sch, name, childPkValues)
		if err != nil {
			return nil, errors.Wrap(err, "RetrieveWithChildren/RetrieveObject("+name+")")
		}
		if obj.Children[name] == nil {
			obj.Children[name] = make(object.ObjectArray, 1)
		}
		obj.Children[name][0] = childObj
	}
	return obj, nil
}

// RetrieveObject function will fleshen an object structure, given some primary keys
func RetrieveObject(ctx context.Context, db *sql.DB, sch *schema.Schema, table string, queryVals map[string]interface{}) (*object.Object, error) {
	// TODO: Implement limit... That's all a singular retrieve should be underneath the hood that's different.
	objAry, err := RetrieveObjects(ctx, db, sch, table, queryVals)
	if err != nil {
		return nil, err
	}
	return objAry[0], nil
}

// RetrieveObjects function will fleshen an object structure, given some primary keys
func RetrieveObjects(ctx context.Context, db *sql.DB, sch *schema.Schema, table string, queryVals map[string]interface{}) (object.ObjectArray, error) {
	objTable := sch.Tables[table]
	if objTable == nil {
		return nil, errors.New("RetrieveObjects: unknown object table " + table)
	}
	var objectArray object.ObjectArray

	gen := sqlitegen.New("sqlite", "test", sch)

	queryObj := object.New(table)
	queryObj.KV = queryVals

	sqlStr, bindArgs, err := gen.BindingRetrieve(sch, queryObj)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(sqlStr)
	//	fmt.Println(bindArgs)

	stmt, err := db.PrepareContext(ctx, sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.QueryContext(ctx, bindArgs...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	columnNames, err := res.Columns()
	if err != nil {
		return nil, err
	}
	columnTypes, err := res.ColumnTypes()
	if err != nil {
		return nil, err
	}

	for res.Next() {
		columnPointers := make([]interface{}, len(columnNames))

		for i := 0; i < len(columnNames); i++ {
			ct := columnTypes[i]
			// TODO: Improve database type support.
			//fmt.Println(ct.DatabaseTypeName())

			// TODO: Do I need to reset columnPointers every time?
			if ct.DatabaseTypeName() == "text" {
				nullable, _ := ct.Nullable()
				if nullable {
					var s sql.NullString
					columnPointers[i] = &s
				} else {
					var s string
					columnPointers[i] = &s
				}
			} else {
				nullable, _ := ct.Nullable()
				if nullable {
					var j sql.NullInt64
					columnPointers[i] = &j
				} else {
					var j int64
					columnPointers[i] = &j

				}
			}
			// columnPointers[i] = &columns[i]
		}

		obj := object.New(table)
		if err := res.Scan(columnPointers...); err != nil {
			return nil, err
		}
		for i, v := range columnPointers {
			ct := columnTypes[i]
			// TODO: Improve database type support.
			//fmt.Println("dbtypename=", ct.DatabaseTypeName())

			if ct.DatabaseTypeName() == "text" {
				nullable, _ := ct.Nullable()
				if nullable {
					val := v.(*sql.NullString)
					if val.Valid {
						obj.Set(columnNames[i], val.String)
					}
					// TODO: We don't set keys for null values. How else can we support this?
				} else {
					val := v.(*string)
					obj.Set(columnNames[i], *val)

				}
			} else {
				// TODO: support more than 'int64' for integer...
				nullable, _ := ct.Nullable()
				if nullable {
					val := v.(*sql.NullInt64)
					if val.Valid {
						obj.Set(columnNames[i], val.Int64)
					}
					// TODO: We don't set keys for null values. How else can we support this?
				} else {
					val := v.(*int64)
					obj.Set(columnNames[i], *val)
				}
			}
		}
		obj.SetSaved(true)
		obj.ResetChangedFields()

		objectArray = append(objectArray, obj)
	}

	err = res.Err()
	if err != nil {
		return nil, err
	}
	return objectArray, nil
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
