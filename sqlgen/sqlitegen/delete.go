package sqlitegen

import (
	"errors"
	"fmt"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// BindingDelete generates the SQL for a given UPDATE statement for SQLite with binding parameter values
func (g Generator) BindingDelete(sch *schema.Schema, obj *object.Object) (string, []interface{}, error) {
	schTable, ok := sch.Tables[obj.Type]
	if !ok {
		return "", nil, errors.New("BindingDelete: Table map unavailable for table " + obj.Type)
	}

	fieldsMap := schTable.Fields
	if fieldsMap == nil {
		return "", nil, errors.New("BindingDelete: Field map unavailable for table " + obj.Type)
	}

	whereClause, bindWhere, err := renderUpdateWhereClause(schTable, fieldsMap, obj)
	if err != nil {
		return "", nil, err
	}

	// TODO: use schema name from object lookup type, fix in other places...
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE %s", obj.Type, whereClause)
	//fmt.Println(sqlStr)
	return sqlStr, bindWhere, nil
}
