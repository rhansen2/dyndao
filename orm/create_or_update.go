package orm

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/object"
)

func (o ORM) CreateOrUpdateTx(ctx context.Context, tx *sql.Tx, obj *object.Object) (int64, *object.Object, error) {
	var err error
	var retObj *object.Object

	retObj, err = o.RetrieveTx(ctx, tx, obj.Type, obj.KV)
	if err != nil {
		return 0, nil, err
	}

	var numRows int64
	var opType string

	if retObj == nil {
		numRows, err = o.Insert(ctx, tx, obj)
		opType = "Insert"
	} else {
		numRows, err = o.Update(ctx, tx, obj)
		opType = "Update"
	}

	if err != nil {
		return 0, nil, err
	}
	if numRows == 0 {
		return 0, nil, errors.New("CreateOrUpdateTx: numRows was 0 when expecting " + opType)
	}

	return numRows, obj, nil
}

func (o ORM) CreateOrUpdate(ctx context.Context, obj *object.Object) (int64, *object.Object, error) {
	return o.CreateOrUpdateTx(ctx, nil, obj)
}

func (o ORM) CreateOrUpdateKVTx(ctx context.Context, tx *sql.Tx, typ string, queryKV map[string]interface{}, createKV map[string]interface{}) (int64, *object.Object, error) {
	var err error
	var retObj *object.Object

	retObj, err = o.RetrieveTx(ctx, tx, typ, queryKV)
	if err != nil {
		return 0, nil, err
	}

	var numRows int64
	var opType string

	var obj *object.Object
	if retObj == nil {
		obj = object.New(typ)
		obj.KV = createKV
		numRows, err = o.Insert(ctx, tx, obj)
		opType = "Insert"
	} else {
		obj = object.New(typ)
		obj.KV = createKV
		numRows, err = o.Update(ctx, tx, obj)
		opType = "Update"
	}

	if err != nil {
		return 0, nil, err
	}
	if numRows == 0 {
		return 0, nil, errors.New("CreateOrUpdateKVTx: numRows was 0 when expecting " + opType)
	}

	return numRows, obj, nil
}

func (o ORM) CreateOrUpdateKV(ctx context.Context, typ string, queryKV map[string]interface{}, createKV map[string]interface{}) (int64, *object.Object, error) {
	return o.CreateOrUpdateKVTx(ctx, nil, typ, queryKV, createKV)
}
