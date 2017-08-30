package oraclegen

import (
	"os"
	"errors"
	"fmt"
	"strings"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// BindingUpdate generates the SQL for a given UPDATE statement for oracle with binding parameter values
func (g Generator) BindingUpdate(sch *schema.Schema, obj *object.Object) (string, []interface{}, []interface{}, error) {
	schTbl := sch.GetTable(obj.Type)
	if schTbl == nil {
		return "", nil, nil, errors.New("BindingUpdate: Table map unavailable for table " + obj.Type)
	}

	fieldsMap := schTbl.Fields
	if fieldsMap == nil {
		return "", nil, nil, errors.New("BindingUpdate: Field map unavailable for table " + obj.Type)
	}

	whereClause, bindWhere, err := renderUpdateWhereClause(schTbl, fieldsMap, obj)
	if err != nil {
		return "", nil, nil, err
	}

	i := 0

	var bindArgs []interface{}
	var newValuesAry []string

	// If some things have changed, then only use fields that we're sure have changed
	if len(obj.ChangedFields) > 0 {
		bindArgs = make([]interface{}, len(obj.ChangedFields))
		newValuesAry = make([]string, len(obj.ChangedFields))

		for k := range obj.ChangedFields {
			f := schTbl.GetField(k)
			if f.IsIdentity {
				continue
			}
			v := obj.KV[k]
			if v == nil {
				newValuesAry[i] = fmt.Sprintf("%s = NULL", f.Name)
				bindArgs[i] = nil
			} else {
				newValuesAry[i] = fmt.Sprintf("%s = %s%d", f.Name, renderBindingUpdateValue(f), i)
				bindArgs[i] = v
			}
			i++
		}
	} else {
		// An update where it's not explicitly clear that anything has changed should
		// just set every field we have available.
		bindArgs = make([]interface{}, len(obj.KV)-1)
		// TODO: -1 for Oracle because we expect an identity field
		newValuesAry = make([]string, len(obj.KV)-1)

		for k, v := range obj.KV {
			f := schTbl.GetField(k)
			if f.IsIdentity {
				continue
			}
			if v == nil {
				newValuesAry[i] = fmt.Sprintf("%s = NULL", f.Name)
				bindArgs[i] = nil
			} else {
				newValuesAry[i] = fmt.Sprintf("%s = %s%d", f.Name, renderBindingUpdateValue(f), i)
				bindArgs[i] = v
			}

			i++
		}
	}
	bindArgs = removeNilsIfNeeded(bindArgs)

	tableName := schema.GetTableName(schTbl.Name, obj.Type)
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, strings.Join(newValuesAry, ","), whereClause)
	if os.Getenv("DEBUG_UPDATE") != "" {
		fmt.Println("BindingUpdate/sqlStr->",sqlStr)
	}
	return sqlStr, bindArgs, bindWhere, nil
}

func renderBindingUpdateValue(f *schema.Field) string {
	return ":" + f.Name
}
