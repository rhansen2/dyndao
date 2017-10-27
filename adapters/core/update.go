package core

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
	"github.com/rbastic/nils"
)

func zeroTime(arg interface{}) bool {
	switch t := arg.(type) {
	case time.Time:
		if t.IsZero() {
			return true
		}
	case *time.Time:
		if t == nil {
			return true
		}
		if t.IsZero() {
			return true
		}
	}
	return false
}

func sqlValueConvert(arg interface{}) (string, bool) {
	switch arg.(type) {
	case *object.SQLValue:
		v := arg.(*object.SQLValue)
		return v.String(), true
	default:
		return "", false
	}
}

func safeConvert(arg interface{}) time.Time {
	switch t := arg.(type) {
	case string:
		tt, err := time.Parse(time.RFC3339, t)
		if err != nil {
			panic(err)
		}
		return tt
	case time.Time:
		return t
	case *time.Time:
		return *t
	default:
		panic(fmt.Sprintf("unknown type in safe convert: %v", reflect.TypeOf(t)))
	}
}

// BindingUpdate generates the SQL for a given UPDATE statement for oracle with binding parameter values
func BindingUpdate(g *sg.SQLGenerator, sch *schema.Schema, obj *object.Object) (string, []interface{}, []interface{}, error) {
	schTbl := sch.GetTable(obj.Type)
	if schTbl == nil {
		return "", nil, nil, errors.New("BindingUpdate: Table map unavailable for table " + obj.Type)
	}

	fieldsMap := schTbl.Columns
	if fieldsMap == nil {
		return "", nil, nil, errors.New("BindingUpdate: Column map unavailable for table " + obj.Type)
	}

	whereClause, bindWhere, err := g.RenderUpdateWhereClause(g, schTbl, fieldsMap, obj)
	if err != nil {
		return "", nil, nil, err
	}

	var bindArgs []interface{}
	var newValuesAry []string
	var kv map[string]interface{}

	if len(obj.ChangedColumns) > 0 {
		bindArgs = make([]interface{}, len(obj.ChangedColumns))
		newValuesAry = make([]string, len(obj.ChangedColumns))
		kv = obj.ChangedColumns
	} else {
		// An update where it's not explicitly clear that anything has changed should
		// just set every field we have available.
		bindArgs = make([]interface{}, len(obj.KV)-1)
		// TODO: -1 for Oracle because we expect an identity field
		// TODO FIXME: Is this correct / incorrect still?
		newValuesAry = make([]string, len(obj.KV)-1)
		kv = obj.KV
	}

	i := 0
	for k, v := range kv {
		f := schTbl.GetColumn(k)
		if f == nil {
			return "", nil, nil, errors.New("BindingUpdate: field config unavailable for object Type: " + obj.Type + ", key: " + k)
		}
		if f.IsIdentity {
			continue
		}

		vStr, wasSV := sqlValueConvert(v)
		if wasSV {
			newValuesAry[i] = fmt.Sprintf("%s = %s", f.Name, vStr)
			bindArgs[i] = nil
		} else {
			if g.IsTimestampType(schTbl.GetColumn(k).DBType) {
				v = safeConvert(v)
			}
			if v == nil || zeroTime(v) {
				newValuesAry[i] = fmt.Sprintf("%s = NULL", f.Name)
				bindArgs[i] = nil
			} else {
				newValuesAry[i] = fmt.Sprintf("%s = %s", f.Name, g.RenderBindingValueWithInt(f, int64(i)))
				bindArgs[i] = v
			}
		}
		i++
	}

	bindArgs = nils.RemoveNilsIfNeeded(bindArgs)

	tableName := schema.GetTableName(schTbl.Name, obj.Type)
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, strings.Join(newValuesAry, ","), whereClause)
	return sqlStr, bindArgs, bindWhere, nil
}
