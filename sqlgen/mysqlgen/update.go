package mysqlgen

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
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
		panic("unkown type in safe convert")
	}
}

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

	whereClause, bindWhere, err := g.renderUpdateWhereClause(schTbl, fieldsMap, obj)
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
			if g.IsTimestampType(schTbl.GetField(k).DBType) {
				v = safeConvert(v)
			}
			if v == nil || zeroTime(v) {
				newValuesAry[i] = fmt.Sprintf("%s = NULL", f.Name)
				bindArgs[i] = nil
			} else {
				newValuesAry[i] = fmt.Sprintf("%s = %s", f.Name, g.RenderBindingValue(f))
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
			if g.IsTimestampType(schTbl.GetField(k).DBType) {
				v = safeConvert(v)
			}
			if v == nil || zeroTime(v) {
				newValuesAry[i] = fmt.Sprintf("%s = NULL", f.Name)
				bindArgs[i] = nil
			} else {
				newValuesAry[i] = fmt.Sprintf("%s = %s", f.Name, g.RenderBindingValue(f))
				bindArgs[i] = v
			}

			i++
		}
	}
	bindArgs = removeNilsIfNeeded(bindArgs)

	tableName := schema.GetTableName(schTbl.Name, obj.Type)
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, strings.Join(newValuesAry, ","), whereClause)
	return sqlStr, bindArgs, bindWhere, nil
}

