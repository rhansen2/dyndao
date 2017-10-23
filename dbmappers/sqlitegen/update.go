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
	schTbl := sch.GetTable(obj.Type)
	if schTbl == nil {
		return "", nil, nil, errors.New("BindingUpdate: Table map unavailable for table " + obj.Type)
	}

	fieldsMap := schTbl.Fields
	if fieldsMap == nil {
		return "", nil, nil, errors.New("BindingUpdate: Field map unavailable for table " + obj.Type)
	}

	whereClause, bindWhere, err := g.renderUpdateWhereClause(schTbl, fieldsMap, obj)
	if err != nil {
		return "", nil, nil, err
	}

	bindArgs := make([]interface{}, len(obj.KV))
	newValuesAry := make([]string, len(obj.KV))
	var i int64
	for k, v := range obj.KV {
		f := fieldsMap[k]
		newValuesAry[i] = fmt.Sprintf("%s = %s", f.Name, g.RenderBindingValueWithInt(f, i))
		bindArgs[i] = v
		i++
	}

	tableName := schema.GetTableName(schTbl.Name, obj.Type)
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, strings.Join(newValuesAry, ","), whereClause)
	return sqlStr, bindArgs, bindWhere, nil
}
