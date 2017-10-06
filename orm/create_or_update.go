package orm

import (
	"context"
	"database/sql"
	"github.com/rbastic/dyndao/object"
	"github.com/pkg/errors"
)

func (o ORM) CreateOrUpdateTx(ctx context.Context, tx * sql.Tx, obj * object.Object) (int, *object.Object, error) {
	retObj, err := o.RetrieveTx(ctx, tx, obj.Type, obj.KV)
	if err != nil {
		return 0, nil, err
	}

	if retObj == nil {
		numRows, err := o.Insert(ctx, tx, obj)
		if err != nil {
			return 0, nil, err
		}
		if numRows == 0 {
			return 0, nil, errors.New("FindOrCreate: numRows was 0 when expecting Insert")
		}

		return numRows, obj, nil
	} else {
		numRows, err := o.Update(ctx, tx, obj)
		if err != nil {
			return 0, nil, err
		}
		if numRows == 0 {
			return 0, nil, errors.New("FindOrCreate: numRows was 0 when expecting Update")
		}

		return numRows, obj, nil
	}

	return 0, nil, nil
}

func (o ORM) CreateOrUpdate(ctx context.Context, table string, pkVals map[string]interface{}, createVals map[string]interface{}) (int, *object.Object, error) {

	return o.CreateOrUpdate(ctx, nil, table, queryVals, createVals)
}
