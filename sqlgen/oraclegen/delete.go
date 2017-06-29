package oraclegen

import (
	"errors"
	"fmt"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// BindingDelete generates the SQL for a given UPDATE statement for SQLite with binding parameter values
func (g Generator) BindingDelete(sch *schema.Schema, obj * object.Object) (string, []interface{}, []interface{}, error) {
	schTable, ok := sch.Tables[obj.Type]
	if !ok {
		return "", nil, nil, errors.New("BindingDelete: Table map unavailable for table " + obj.Type)
	}

	fieldsMap := schTable.Fields
	if fieldsMap == nil {
		return "", nil, nil, errors.New("BindingDelete: Field map unavailable for table " + obj.Type)
	}

	whereClause, bindWhere, err := renderUpdateWhereClause(schTable, fieldsMap, obj)
	if err != nil {
		return "", nil, nil, err
	}

	bindArgs := make([]interface{}, len(obj.KV))
	i := 0
	for _, v := range obj.KV {
		bindArgs[i] = v
		i++
	}

	// TODO: use schema name from object lookup type, fix in other places...
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE %s", obj.Type, whereClause)
	//fmt.Println(sqlStr)
	return sqlStr, bindArgs, bindWhere, nil
}
