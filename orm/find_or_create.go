package orm

import (
	"context"
	"database/sql"
	"github.com/rbastic/dyndao/object"
	"github.com/pkg/errors"
)

func (o ORM) FindOrCreate(ctx context.Context, tx * sql.Tx, table string, queryVals map[string]interface{}) (*object.Object, error) {
	obj, err := o.RetrieveTx(ctx, tx, table, queryVals)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		obj = object.New( table )
		obj.KV = queryVals
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
