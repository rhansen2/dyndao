package orm

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/object"
)

func (o ORM) FindOrCreateTx(ctx context.Context, tx *sql.Tx, obj *object.Object) (*object.Object, error) {
	obj, err := o.RetrieveTx(ctx, tx, obj.Type, obj.KV)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		numRows, err := o.Insert(ctx, tx, obj)
		if err != nil {
			return nil, err
		}
		if numRows == 0 {
			return nil, errors.New("FindOrCreate: numRows was 0 when expecting Insert")
		}

		return obj, nil
	}

	return obj, nil
}

func (o ORM) FindOrCreate(ctx context.Context, obj *object.Object) (*object.Object, error) {

	return o.FindOrCreateTx(ctx, nil, obj)
}
