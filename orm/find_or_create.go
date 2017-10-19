package orm

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/object"
)

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

func (o ORM) FindOrCreate(ctx context.Context, obj *object.Object) (int64, *object.Object, error) {
	return o.FindOrCreateTx(ctx, nil, obj)
}

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

func (o ORM) FindOrCreateKV(ctx context.Context, typ string, queryKV map[string]interface{}, createKV map[string]interface{}) (int64, *object.Object, error) {
	return o.FindOrCreateKVTx(ctx, nil, typ, queryKV, createKV)
}
