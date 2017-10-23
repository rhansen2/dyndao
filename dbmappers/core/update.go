package core

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	sg "github.com/rbastic/dyndao/sqlgen"
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
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
func BindingUpdate(g * sg.SQLGenerator, sch *schema.Schema, obj *object.Object) (string, []interface{}, []interface{}, error) {
	schTbl := sch.GetTable(obj.Type)
	if schTbl == nil {
		return "", nil, nil, errors.New("BindingUpdate: Table map unavailable for table " + obj.Type)
	}

	fieldsMap := schTbl.Fields
	if fieldsMap == nil {
		return "", nil, nil, errors.New("BindingUpdate: Field map unavailable for table " + obj.Type)
	}

	whereClause, bindWhere, err := g.RenderUpdateWhereClause(g, schTbl, fieldsMap, obj)
	if err != nil {
		return "", nil, nil, err
	}

	i := 0

	var bindArgs []interface{}
	var newValuesAry []string

	// If some things have changed, then only use fields that we're sure have changed

	// TODO: Refactor this code.
	if len(obj.ChangedFields) > 0 {
		bindArgs = make([]interface{}, len(obj.ChangedFields))
		newValuesAry = make([]string, len(obj.ChangedFields))

		for k := range obj.ChangedFields {
			f := schTbl.GetField(k)
			if f == nil {
				return "", nil, nil, errors.New("BindingUpdate: field config unavailable for object Type: " + obj.Type + ", key: " + k)
			}
			if f.IsIdentity {
				continue
			}
			v := obj.KV[k]

			vStr, wasSV := sqlValueConvert(v)
			if wasSV {
				newValuesAry[i] = fmt.Sprintf("%s = %s", f.Name, vStr)
				bindArgs[i] = nil
			} else {
				if g.IsTimestampType(schTbl.GetField(k).DBType) {
					v = safeConvert(v)
				}
				if v == nil || zeroTime(v) {
					newValuesAry[i] = fmt.Sprintf("%s = NULL", f.Name)
					bindArgs[i] = nil
				} else {
					newValuesAry[i] = fmt.Sprintf("%s = %s%d", f.Name, g.RenderBindingValue(f), i)
					bindArgs[i] = v
				}
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
				if g.IsTimestampType(schTbl.GetField(k).DBType) {
					v = safeConvert(v)
				}
				if v == nil || zeroTime(v) {
					newValuesAry[i] = fmt.Sprintf("%s = NULL", f.Name)
					bindArgs[i] = nil
				} else {
					newValuesAry[i] = fmt.Sprintf("%s = %s%d", f.Name, g.RenderBindingValue(f), i)
					bindArgs[i] = v
				}
			}

			i++
		}
	}
	bindArgs = nils.RemoveNilsIfNeeded(bindArgs)

	tableName := schema.GetTableName(schTbl.Name, obj.Type)
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, strings.Join(newValuesAry, ","), whereClause)
	return sqlStr, bindArgs, bindWhere, nil
}
