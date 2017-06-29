package orm

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

func pkQueryValsFromKV(obj *object.Object, sch *schema.Schema, parentTableName string) (map[string]interface{}, error) {
	qv := make(map[string]interface{})

	schemaTable := sch.Tables[parentTableName]
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

// GetParentsViaChild retrieves all direct (one-level 'up') parents for a given child object.
// If a child contains multiple parent tables (possibility?) then this would return an Array
// of objects with multiple potential values for their obj.Type fields.
func (o ORM) GetParentsViaChild(ctx context.Context, childObj *object.Object) (object.Array, error) {
	table := childObj.Type

	objTable := o.s.Tables[table]
	if objTable == nil {
		return nil, errors.New("GetParentsViaChild: unknown object table " + table)
	}

	var parentObjs object.Array

	if objTable.ParentTables == nil {
		return nil, errors.New("GetParentsViaChild: cannot retrieve parents for table " + table + ", schema ParentTables is nil")
	}
	for _, pt := range objTable.ParentTables {
		pkQueryVals, err := pkQueryValsFromKV(childObj, o.s, pt)
		if err != nil {
			return nil, err
		}
		objs, err := o.RetrieveObjects(ctx, pt, pkQueryVals)
		if err != nil {
			return nil, err
		}
		parentObjs = append(parentObjs, objs...)
	}

	return parentObjs, nil
}

// TODO: For foreign key filling, we do not check to see if there are conflicts
// with regards to the uniqueness of primary key names.

// RetrieveWithChildren function will fleshen an *entire* object structure, given some primary keys
func (o ORM) RetrieveWithChildren(ctx context.Context, table string, pkValues map[string]interface{}) (*object.Object, error) {
	objTable := o.s.Tables[table]
	if objTable == nil {
		return nil, errors.New("RetrieveWithChildren: unknown object table " + table)
	}

	obj, err := o.RetrieveObject(ctx, table, pkValues)
	if err != nil {
		return nil, errors.Wrap(err, "RetrieveWithChildren/RetrieveObject")
	}

	for name := range objTable.Children {
		childPkValues := make(map[string]interface{})

		childSchemaTable := o.s.Tables[name]

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
		childObj, err := o.RetrieveObject(ctx, name, childPkValues)
		if err != nil {
			return nil, errors.Wrap(err, "RetrieveWithChildren/RetrieveObject("+name+")")
		}
		if obj.Children[name] == nil {
			obj.Children[name] = make(object.Array, 1)
		}
		obj.Children[name][0] = childObj
	}
	return obj, nil
}

// RetrieveObject function will fleshen an object structure, given some primary keys
func (o ORM) RetrieveObject(ctx context.Context, table string, queryVals map[string]interface{}) (*object.Object, error) {
	// TODO: Implement LIMIT... That's all a singular retrieve should be underneath the hood that's different.
	objAry, err := o.RetrieveObjects(ctx, table, queryVals)
	if err != nil {
		return nil, err
	}
	if objAry == nil {
		return nil, errors.New("RetrieveObject: received empty object array from RetrieveObjects")
	}
	return objAry[0], nil
}

// FleshenChildren function accepts an object and resets it's children.
func (o ORM) FleshenChildren(ctx context.Context, table string, obj *object.Object) (*object.Object, error) {
	schemaTable := o.s.Tables[obj.Type]
	pkKey := schemaTable.Primary
	pkVal := obj.Get(pkKey)

	if len(schemaTable.Children) > 0 {
		// Does this table have child tables?
		for childTableName := range schemaTable.Children {
			m := map[string]interface{}{}
			m[pkKey] = pkVal
			childObjs, err := o.RetrieveObjects(ctx, childTableName, m)
			if err != nil {
				return nil, err
			}
			obj.Children[childTableName] = childObjs
		}
	}
	return obj, nil
}

// RetrieveObjects function will fleshen an object structure, given some primary keys
func (o ORM) RetrieveObjects(ctx context.Context, table string, queryVals map[string]interface{}) (object.Array, error) {
	if o.s == nil {
		return nil, errors.New("RetrieveObjects: why is ORM schema set to nil?")
	}
	if o.s.Tables == nil {
		return nil, errors.New("RetrieveObjects: why is Tables nil?")
	}
	objTable := o.s.Tables[table]
	if objTable == nil {
		return nil, errors.New("RetrieveObjects: unknown object table " + table)
	}
	var objectArray object.Array

	queryObj := object.New(table)
	queryObj.KV = queryVals

	sqlStr, columnNames, bindArgs, err := o.sqlGen.BindingRetrieve(o.s, queryObj)
	if os.Getenv("DEBUG") != "" {
		fmt.Println("RetrieveObjects/sqlStr=", sqlStr, "columnNames=", columnNames, "bindArgs=", bindArgs)
	}

	if err != nil {
		return nil, err
	}

	stmt, err := o.RawConn.PrepareContext(ctx, sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.QueryContext(ctx, bindArgs...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	columnTypes, err := res.ColumnTypes()
	if err != nil {
		return nil, err
	}

	for res.Next() {
		columnPointers, err := o.makeColumnPointers(len(columnNames), columnTypes)
		if err != nil {
			return nil, err
		}

		obj := object.New(table)
		if err := res.Scan(columnPointers...); err != nil {
			return nil, err
		}

		err = o.dynamicObjectSetter(columnNames, columnPointers, columnTypes, obj)
		if err != nil {
			return nil, err
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

func (o ORM) dynamicObjectSetter(columnNames []string, columnPointers []interface{}, columnTypes []*sql.ColumnType, obj *object.Object) error {
	sqlGen := o.sqlGen
	for i, v := range columnPointers {
		ct := columnTypes[i]

		typeName := ct.DatabaseTypeName()
		if sqlGen.IsStringType(typeName) {
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
		} else if sqlGen.IsNumberType(typeName) {
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
		} else {
			return errors.New("dynamicObjectSetter: Unrecognized type: " + typeName)
		}
	}
	return nil
}

func (o ORM) makeColumnPointers(sliceLen int, columnTypes []*sql.ColumnType) ([]interface{}, error) {
	columnPointers := make([]interface{}, sliceLen)
	sqlGen := o.sqlGen
	for i := 0; i < sliceLen; i++ {
		ct := columnTypes[i]
		typeName := ct.DatabaseTypeName()
		if sqlGen.IsStringType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				var s sql.NullString
				columnPointers[i] = &s
			} else {
				var s string
				columnPointers[i] = &s
			}
		} else if sqlGen.IsNumberType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				var j sql.NullInt64
				columnPointers[i] = &j
			} else {
				var j int64
				columnPointers[i] = &j

			}
		} else {
			return nil, errors.New("makeColumnPointers: Unrecognized type: " + typeName)
		}
	}
	return columnPointers, nil
}

// TODO: Read this post for more info on the above...
// https://stackoverflow.com/questions/23507531/is-golangs-sql-package-incapable-of-ad-hoc-exploratory-queries/23507765#23507765
