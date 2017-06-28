package oraclegen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// BindingUpdate generates the SQL for a given UPDATE statement for SQLite with binding parameter values
func (g Generator) BindingUpdate(sch *schema.Schema, obj *object.Object) (string, []interface{}, []interface{}, error) {
	schTable, ok := sch.Tables[obj.Type]
	if !ok {
		return "", nil, nil, errors.New("BindingUpdate: Table map unavailable for table " + obj.Type)
	}

	fieldsMap := schTable.Fields
	if fieldsMap == nil {
		return "", nil, nil, errors.New("BindingUpdate: Field map unavailable for table " + obj.Type)
	}

	whereClause, bindWhere, err := renderUpdateWhereClause(schTable, fieldsMap, obj)
	if err != nil {
		return "", nil, nil, err
	}

	bindArgs := make([]interface{}, len(obj.KV)-1)
	// TODO: - 1 for Oracle because we expect an identity field
	newValuesAry := make([]string, len(obj.KV)-1)
	i := 0
	for k, v := range obj.KV {
		f := fieldsMap[k]

		if f.IsIdentity {
			continue
		}

		newValuesAry[i] = fmt.Sprintf("%s = %s%d", f.Name, renderBindingUpdateValue(f), i)
		bindArgs[i] = v
		fmt.Println("newValues[i]=", newValuesAry[i], "bindArgs[i]=", bindArgs[i])

		i++
	}

	// TODO: use schema name from object lookup type, fix in other places...
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s", obj.Type, strings.Join(newValuesAry, ","), whereClause)
	fmt.Println(sqlStr)
	fmt.Println(bindArgs)
	fmt.Println(bindWhere)
	return sqlStr, bindArgs, bindWhere, nil
}

func renderBindingUpdateValue(f *schema.Field) string {
	return ":" + f.Name
}