package orm

// The ORM package is designed to tie everything together: a database connection, schema,
// relevant objects, etc. The current design is a WIP. While not finished, it is serviceable
// and can be used effectively.

/*

TODO: This file has two Retrieve implementations with some duplicated
code (RetrieveManyFromCustomSQL, retrieveManyCore). It should be easy to
refactor.

*/

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// GetParentsViaChild retrieves all direct (one-level 'up') parents for a given child object.
// If a child contains multiple parent tables (possibility?) then this would return an Array
// of objects with multiple potential values for their obj.Type fields.
func (o ORM) GetParentsViaChild(ctx context.Context, childObj *object.Object) (object.Array, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	table := childObj.Type

	// Retrieve schema table object
	objTable := o.s.GetTable(table)
	if objTable == nil {
		return nil, errors.New("GetParentsViaChild: unknown object table " + table)
	}

	var parentObjs object.Array
	// This method can only work if the schema is configured hierarchically
	if objTable.ParentTables == nil {
		return nil, errors.New("GetParentsViaChild: cannot retrieve parents for table " + table + ", schema ParentTables is nil")
	}
	for _, pt := range objTable.ParentTables {
		// 'Capture' the primary key(s)
		pkQueryVals, err := pkQueryValsFromKV(childObj, o.s, pt)
		if err != nil {
			return nil, err
		}
		// Retrieve + append the relevant parent objs
		objs, err := o.RetrieveMany(ctx, pt, pkQueryVals)
		if err != nil {
			return nil, err
		}
		parentObjs = append(parentObjs, objs...)
	}

	return parentObjs, nil
}

// RetrieveWithChildren function will fleshen an *entire* object structure, given some primary keys
// By entire, we mean that we also retrieve any relevant children objects. However, we do not call RetrieveWithChildren
// when fleshening the children structures -- when retrieving the children, we do a single-level retrieve, ignoring
// any child structures that may be configured at two levels of depth.
func (o ORM) RetrieveWithChildren(ctx context.Context, table string, pkValues map[string]interface{}) (*object.Object, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	// Retrieve schema.Table object
	objTable := o.s.GetTable(table)
	if objTable == nil {
		return nil, errors.New("RetrieveWithChildren: unknown object table " + table)
	}

	// Retrieve single object from database
	obj, err := o.Retrieve(ctx, table, pkValues)
	if err != nil {
		return nil, errors.Wrap(err, "RetrieveWithChildren/Retrieve")
	}

	// Iterate the configured 'children' for this particular object type
	for name := range objTable.Children {
		childPkValues := make(map[string]interface{})

		// Retrieve the active 'schema table' for this child
		childSchemaTable := o.s.GetTable(name)
		if childSchemaTable == nil {
			return nil, fmt.Errorf("RetrieveWithChildren: unknown object table for child type %s", name)
		}

		// (For each child...) Propagate the 'primary key value' from the parent object if needed.
		childPkName := childSchemaTable.Primary
		pVal, ok := pkValues[childPkName]
		if ok {
			childPkValues[childPkName] = pVal
		}

		// Propagate foreign key values for retrieval
		if childSchemaTable.MultiKey && childSchemaTable.ForeignKeys != nil {
			for _, fk := range childSchemaTable.ForeignKeys {
				// TODO: Check that value exists before we
				// attempt to set?
				/*
					v, ok := pkValues[fk]
					if ok {
						childPkValues[fk] = v
					} else {
						// TODO: is this an error condition?
					}
				*/

				childPkValues[fk] = pkValues[fk]
			}
		}

		// Retrieve a single child (TODO: Implement RetrieveMany options as well)
		childObj, err := o.Retrieve(ctx, name, childPkValues)
		if err != nil {
			return nil, errors.Wrap(err, "RetrieveWithChildren/Retrieve("+name+")")
		}
		// Populate single child inside parent
		if obj.Children[name] == nil {
			obj.Children[name] = make(object.Array, 1)
		}
		obj.Children[name][0] = childObj
	}

	return obj, nil
}

// retrieveCore function will fleshen an object structure, given some primary keys.
// Technically, we call RetrieveMany internally. Since we do not have LIMIT implemented yet,
// it's just a cheap implementation that returns the zeroeth value. Nil will be returned
// for both the object and the error if a row is unable to be matched by the underlying
// datastore.
// TODO: Implement LIMIT so that we can improve this.
func (o ORM) retrieveCore(ctx context.Context, tx *sql.Tx, table string, queryVals map[string]interface{}) (*object.Object, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	objAry, err := o.retrieveManyCore(ctx, tx, table, queryVals)
	if err != nil {
		return nil, err
	}
	// We return nil, nil to indicate a lack of result.
	if objAry == nil {
		return nil, nil
	}
	return objAry[0], nil
}

// RetrieveTx function will fleshen an object structure, given some primary keys.
// Technically, we call RetrieveMany internally. Since we do not have LIMIT implemented yet,
// it's just a cheap implementation that returns the zeroeth value. Nil will be returned
// for both the object and the error if a row is unable to be matched by the underlying
// datastore.
// TODO: Implement LIMIT so that we can improve this.
func (o ORM) RetrieveTx(ctx context.Context, tx *sql.Tx, table string, queryVals map[string]interface{}) (*object.Object, error) {
	return o.retrieveCore(ctx, tx, table, queryVals)
}

// Retrieve function will fleshen an object structure, given some primary keys.
// Technically, we call RetrieveMany internally. Since we do not have LIMIT implemented yet,
// it's just a cheap implementation that returns the zeroeth value. Nil will be returned
// for both the object and the error if a row is unable to be matched by the underlying
// datastore.
// TODO: Implement LIMIT so that we can improve this.
func (o ORM) Retrieve(ctx context.Context, table string, queryVals map[string]interface{}) (*object.Object, error) {
	return o.retrieveCore(ctx, nil, table, queryVals)
}

// FleshenChildren function accepts an object and resets it's children.
func (o ORM) FleshenChildren(ctx context.Context, obj *object.Object) (*object.Object, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	// Retrieve schema configuration for this object type (schema.Table)
	schemaTable := o.s.GetTable(obj.Type)

	// Retrieve primary key
	pkKey := schemaTable.Primary
	pkVal := obj.Get(pkKey)

	// If this table is configured with child tables, then we iterate over
	// them and call RetrieveMany using the singular primary key value.
	// FIXME: We need to support multikey in this instance if we are going
	// to consider this complete.
	if len(schemaTable.Children) > 0 {
		for childTableName := range schemaTable.Children {
			// TODO: multi-key support here...
			m := map[string]interface{}{}
			m[pkKey] = pkVal
			childObjs, err := o.RetrieveMany(ctx, childTableName, m)
			if err != nil {
				return nil, err
			}
			obj.Children[childTableName] = childObjs
		}
	}
	return obj, nil
}

// RetrieveManyFromCustomSQL will fleshen an object structure, given a custom SQL string. It must still be told
// the column names and the binding arguments in addition to the SQL string, so that it can dynamically map
// the column types accordingly to the destination object. (Mainly, so we know the array length..)
func (o ORM) RetrieveManyFromCustomSQL(ctx context.Context, table string, sqlStr string, columnNames []string, bindArgs []interface{}) (object.Array, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	var objectArray object.Array

	if os.Getenv("DB_TRACE") != "" {
		fmt.Println("RetrieveManyFromCustomSQL/sqlStr=", sqlStr, "columnNames=", columnNames, "bindArgs=", bindArgs)
	}

	stmt, err := o.RawConn.PrepareContext(ctx, sqlStr)
	if err != nil {
		return nil, err
	}
	defer func() {
		stmtErr := stmt.Close()
		if stmtErr != nil {
			fmt.Println("defer stmt.Close error:", stmtErr) // TODO: logger implementation
		}
	}()

	res, err := stmt.QueryContext(ctx, bindArgs...)
	if err != nil {
		return nil, err
	}
	defer func() {
		resErr := res.Close()
		if resErr != nil {
			fmt.Println("defer res.Close error:", resErr) // TODO: logger implementation
		}
	}()

	columnTypes, err := res.ColumnTypes()
	if err != nil {
		return nil, err
	}
	sg := o.sqlGen
	for res.Next() {
		columnPointers, err := sg.MakeColumnPointers(sg, len(columnNames), columnTypes)
		if err != nil {
			return nil, err
		}

		if err := res.Scan(columnPointers...); err != nil {
			return nil, err
		}

		obj := object.New(table)
		err = sg.DynamicObjectSetter(sg, columnNames, columnPointers, columnTypes, obj)
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

func (o ORM) makeQueryObj(objTable *schema.Table, queryVals map[string]interface{}) *object.Object {
	queryObj := object.New(objTable.Name)

	if objTable.FieldAliases == nil {
		queryObj.KV = queryVals
		return queryObj
	}
	for k, v := range queryVals {
		realName := objTable.GetFieldName(k)
		queryObj.KV[realName] = v
	}
	return queryObj
}

// TODO: Have a more powerful querying mechanism, queryVals is very limited.
func (o ORM) retrieveManyCore(ctx context.Context, tx *sql.Tx, table string, queryVals map[string]interface{}) (object.Array, error) {
	// Check for timeout
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Check to make sure that the table arg is a valid table name We will
	// need objTable later.
	objTable := o.s.GetTable(table)
	if objTable == nil {
		return nil, errors.New("RetrieveMany: unknown object table " + table)
	}
	if objTable.Name == "" {
		return nil, errors.New("RetrieveMany: schema table object has unset 'Name' property")
	}

	var objectArray object.Array

	// Construct a dyndao object from our queryVals
	queryObj := o.makeQueryObj(objTable, queryVals)

	// Generate a sql string, the column names, and the binding parameter
	// arguments from the schema and the query object
	sg := o.sqlGen
	sqlStr, columnNames, bindArgs, err := sg.BindingRetrieve(sg, o.s, queryObj)
	if os.Getenv("DB_TRACE") != "" {
		fmt.Println("RetrieveMany/sqlStr=", sqlStr, "columnNames=", columnNames, "bindArgs=", bindArgs)
	}

	if err != nil {
		return nil, err
	}

	// Determines whether we are running inside a transaction or not,
	// returning stmt either way
	stmt, err := stmtFromDbOrTx(ctx, o, tx, sqlStr)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err) // TODO logger implementation
		}
	}()

	res, err := stmt.QueryContext(ctx, bindArgs...)
	if err != nil {
		return nil, err
	}

	defer func() {
		//fmt.Println("RYAN DEFER RETRIEVE ABOUT TO CLOSE")
		err := res.Close()
		if err != nil {
			//fmt.Println("RYAN DEFER RETRIEVE ERROR", err)
			fmt.Println(err) // TODO logger implementation
		}
		//fmt.Println("RYAN DEFER RETRIEVE CLOSED")
	}()

	columnTypes, err := res.ColumnTypes()
	if err != nil {
		return nil, err
	}

	for res.Next() {
		columnPointers, err := sg.MakeColumnPointers(sg, len(columnNames), columnTypes)
		if err != nil {
			return nil, err
		}

		obj := object.New(table)
		if err := res.Scan(columnPointers...); err != nil {
			return nil, err
		}

		err = sg.DynamicObjectSetter(sg, columnNames, columnPointers, columnTypes, obj)
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

// RetrieveManyTx function will fleshen a top-level object structure, given some primary keys. And
// it's transactional!
func (o ORM) RetrieveManyTx(ctx context.Context, tx *sql.Tx, table string, queryVals map[string]interface{}) (object.Array, error) {
	return o.retrieveManyCore(ctx, tx, table, queryVals)
}

// RetrieveMany function will fleshen a top-level object structure, given some primary keys
func (o ORM) RetrieveMany(ctx context.Context, table string, queryVals map[string]interface{}) (object.Array, error) {
	return o.retrieveManyCore(ctx, nil, table, queryVals)
}
