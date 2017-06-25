package sqlitegen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)


// BindingUpdate generates the SQL for a given UPDATE statement for SQLite with binding parameter values
func (g Generator) BindingUpdate(sch *schema.Schema, obj *object.Object) (string, []interface{}, []interface{}, error) {
	schTable, ok := g.Schema.Tables[obj.Type]
	if !ok {
		return "", nil, nil, errors.New("BindingUpdate: Table map unavailable for table " + obj.Type)
	}

	fieldsMap := schTable.Fields
	if fieldsMap == nil {
		return "", nil, nil, errors.New("BindingUpdate: Field map unavailable for table " + obj.Type)
	}

	whereClause, bindWhere, err := renderWhereClause(schTable, fieldsMap, obj)
	if err != nil {
		return "", nil, nil, err
	}

	bindArgs := make([]interface{}, len(obj.KV))
	newValuesAry := make([]string, len(obj.KV))
	i := 0
	for k, v := range obj.KV {
		f := fieldsMap[k]
		newValuesAry[i] = fmt.Sprintf("%s = %s%d", f.Name, renderBindingUpdateValue(f), i)
		bindArgs[i] = v
		i++
	}

	// TODO: use schema name from object lookup type, fix in other places...
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s", obj.Type, strings.Join(newValuesAry, ","), whereClause)
	fmt.Println(sqlStr)
	return sqlStr, bindArgs, bindWhere, nil
}

func renderBindingUpdateValue(f *schema.Field) string {
	return ":" + f.Name
}

