package orm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

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
		return nil, fmt.Errorf("pkQueryValsFromKV: no schema table for table named %s", parentTableName)
	}
	schemaPrimary := schemaTable.Primary

	for fName, field := range schemaTable.Columns {
		if field.IsIdentity || field.IsForeignKey || field.Name == schemaPrimary {
			qv[fName] = obj.Get(fName)
		}
	}
	return qv, nil
}

func (o ORM) recurseAndSave(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	rowsAff, err := o.Save(ctx, tx, obj)
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
				return 0, fmt.Errorf("recurseAndSave: Unknown child object type %s for parent type %s", childObj.Type, obj.Type)
			}
			// TODO: support propagation of additional primary keys
			// that are saved from previous recursive saves ...
			// check if the child schema table contains the
			// parent's primary key field as a name
			_, ok = childTable.Columns[table.Primary]
			if ok {
				// set in the child object if the table contains the primary
				childObj.Set(table.Primary, pkVal)
			}

			// TODO: Likely we can just use pkQueryValsFromKV here,
			// but I would like to add some test cases and verify
			// the functionality under different scenarios.
			aff, err := o.recurseAndSave(ctx, tx, childObj)
			if err != nil {
				return rowsAff + aff, err
			}
		}
	}
	return rowsAff, err
}

// SaveAllInsideTx will attempt to save an entire nested object structure inside of a single transaction.
func (o ORM) SaveAllInsideTx(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	rowsAff, err := o.recurseAndSave(ctx, tx, obj)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return 0, errors.Wrap(err, err2.Error())
		}
		return 0, err
	}
	return rowsAff, nil
}

// SaveAll will attempt to save an entire nested object structure inside of a single transaction.
// It begins the transaction, attempts to recursively save the object and all of it's children,
// and any of the children's children, and then will finally rollback/commit as necessary.
func (o ORM) SaveAll(ctx context.Context, obj *object.Object) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	tx, err := o.RawConn.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	rowsAff, err := o.SaveAllInsideTx(ctx, tx, obj)
	if err != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			return 0, errors.Wrap(err, rollErr.Error())
		}
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			return 0, errors.Wrap(err, rollErr.Error())
		}
		return 0, err
	}
	return rowsAff, nil
}

// SaveButErrorIfUpdate function will INSERT or UPDATE a record. It does not attempt to
// save any of the children. If given a transaction, it will use that to
// attempt to insert the data.
func (o ORM) SaveButErrorIfUpdate(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	objTable := o.s.GetTable(obj.Type)
	// skip if object has invalid type
	if objTable == nil {
		return 0, errors.New("SaveButErrorIfUpdate: unknown object table " + obj.Type)
	}
	// skip if object is clean
	if !obj.IsDirty() {
		return 0, nil
	}
	// retrieve primary key value
	pk := objTable.Primary
	if pk == "" {
		return 0, errors.New("SaveButErrorIfUpdate: empty primary key for " + obj.Type)
	}
	// skip if primary key has no field configuration in table schema
	fieldMap := objTable.Columns
	f := fieldMap[pk]
	if f == nil {
		return 0, errors.New("SaveButErrorIfUpdate: empty field " + pk + " for " + obj.Type)
	}
	// Check the primary key to see if we should insert or update
	_, ok := obj.KV[f.Name]
	if !ok {
		return o.Insert(ctx, tx, obj)
	}
	return 0, fmt.Errorf("SaveButErrorIfUpdate: ORM was told to expect an Insert for this obj: %v", obj)
}

// SaveButErrorIfInsert function will UPDATE a record and error if it
// appears that an INSERT should have been performed. This could be necessary in
// situations where an INSERT would compromise the integrity of the data.  If
// given a transaction, it will use that to attempt to insert the data.
func (o ORM) SaveButErrorIfInsert(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	objTable := o.s.GetTable(obj.Type)
	// skip if object has invalid type
	if objTable == nil {
		return 0, errors.New("SaveButErrorIfInsert: unknown object table " + obj.Type)
	}
	// skip objects that are saved
	if !obj.IsDirty() {
		return 0, nil
	}
	// ensure we have a primary key
	pk := objTable.Primary
	if pk == "" {
		return 0, errors.New("SaveButErrorIfInsert: empty primary key for " + obj.Type)
	}
	// ensure the primary key has a field config
	fieldMap := objTable.Columns
	f := fieldMap[pk]
	if f == nil {
		return 0, errors.New("SaveButErrorIfInsert: empty field " + pk + " for " + obj.Type)
	}
	// Check the primary key to see if we should insert or update
	_, ok := obj.KV[f.Name]
	if !ok {
		return 0, fmt.Errorf("SaveButErrorIfInsert: ORM was told to expect an Update for this obj: %v", obj)
	}
	return o.Update(ctx, tx, obj)
}

// Save function will INSERT or UPDATE a record. It does not attempt to
// save any of the children. If given a transaction, it will use that to
// attempt to insert the data.
func (o ORM) Save(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	objTable := o.s.GetTable(obj.Type)
	// skip if object has invalid type
	if objTable == nil {
		return 0, errors.New("Save: unknown object table " + obj.Type)
	}
	// skip if object is saved
	if !obj.IsDirty() {
		return 0, nil
	}
	// retrieve primary key value
	pk := objTable.Primary
	if pk == "" {
		return 0, errors.New("Save: empty primary key for " + obj.Type)
	}
	// skip if primary key has no field configuration in table schema
	fieldMap := objTable.Columns
	f := fieldMap[pk]
	if f == nil {
		return 0, errors.New("Save: empty field " + pk + " for " + obj.Type)
	}
	// Check the primary key to see if we should insert or update
	_, ok := obj.KV[f.Name]
	if !ok {
		return o.Insert(ctx, tx, obj)
	}
	return o.Update(ctx, tx, obj)
}

// use transaction if needed, otherwise just execute a non-transactionalized operation
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

func maybeDereferenceArgs(arg interface{}) interface{} {
	v := reflect.ValueOf(arg)
	return reflect.Indirect(v).Interface()
}
