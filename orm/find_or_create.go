package orm

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/object"
)

// FindOrCreateTx accepts a context, transaction, and dyndao object. It returns the number of rows affected,
// the resulting dyndao object (either freshly Retrieved or freshly Created) and any error that may have occurred.
// This is one of the recommended methods to use for a FindOrCreate. FindOrCreateKVTx is likely better to use
// in certain situations.
func (o ORM) FindOrCreateTx(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, *object.Object, error) {
	obj, err := o.RetrieveTx(ctx, tx, obj.Type, obj.KV)
	if err != nil {
		return 0, nil, err
	}

	if obj == nil {
		numRows, err := o.Insert(ctx, tx, obj)
		if err != nil {
			return numRows, nil, err
		}
		if numRows == 0 {
			return 0, nil, errors.New("FindOrCreateTx: numRows was 0 when expecting Insert")
		}

		return numRows, obj, nil
	}

	return 0, obj, nil
}

// FindOrCreate is a transactionless FindOrCreate operation. This may not be recommended unless you
// know what you are doing and are comfortable potentially violating transactional integrity.
func (o ORM) FindOrCreate(ctx context.Context, obj *object.Object) (int64, *object.Object, error) {
	return o.FindOrCreateTx(ctx, nil, obj)
}

// FindOrCreateKVTx is useful when you need to execute the Find with specific
// query parameters, but if the record does not exist, then a separate set of
// parameters should be inserted into the database.  This matches a typical
// use-case scenario for a FindOrCreate: primary keys are used to locate the
// row and if the row doesn't exist, then the other values are used for the
// INSERT. Using a transaction is recommended.
func (o ORM) FindOrCreateKVTx(ctx context.Context, tx *sql.Tx, typ string, queryKV map[string]interface{}, createKV map[string]interface{}) (int64, *object.Object, error) {
	obj, err := o.RetrieveTx(ctx, tx, typ, queryKV)
	if err != nil {
		return 0, nil, err
	}

	if obj == nil {
		obj := object.New(typ)
		obj.KV = createKV

		numRows, err := o.Insert(ctx, tx, obj)
		if err != nil {
			return numRows, nil, err
		}
		if numRows == 0 {
			return 0, nil, errors.New("FindOrCreateKVTx: numRows was 0 when expecting Insert")
		}

		return numRows, obj, nil
	}

	return 0, obj, nil
}

// FindOrCreateKV is a transactionless FindOrCreate. This may not be recommended unless you
// know what you are doing and are comfortable potentially violating transactional integrity.
func (o ORM) FindOrCreateKV(ctx context.Context, typ string, queryKV map[string]interface{}, createKV map[string]interface{}) (int64, *object.Object, error) {
	return o.FindOrCreateKVTx(ctx, nil, typ, queryKV, createKV)
}
