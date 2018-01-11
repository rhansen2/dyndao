package core

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

func RenderUpdateWhereClause(g *sg.SQLGenerator, schTable *schema.Table, fieldsMap map[string]*schema.Column, obj *object.Object) (string, []interface{}, *int, error) {
	var emptyInt int
	var bindArgs []interface{}
	var whereClause string

	if len(obj.KV) == 0 {
		return "", nil, &emptyInt, nil
	}

	// TODO: This can likely be refactored a bit.
	bindI := 1
	if !schTable.MultiKey {
		f := fieldsMap[schTable.Primary]
		sqlName := f.Name
		whereClause = fmt.Sprintf("%s = %s", sqlName, g.RenderBindingValueWithInt(f, bindI))
		bindArgs = make([]interface{}, 1)
		bindVal := obj.Get(schTable.Primary)
		if bindVal == nil {
			return "", nil, &emptyInt, errors.New("dyndao: RenderUpdateWhereClause: missing primary key " + schTable.Primary)
		}
		bindArgs[0] = bindVal
		bindI++
	} else {
		// MultiKey means that there could be more than just a single primary key
		// on a table. In this case, we definitely care about involving the entire
		// composite key in the index.
		foreignKeyLen := 0
		if schTable.ForeignKeys != nil {
			foreignKeyLen = len(schTable.ForeignKeys)
		}

		bindArgsLen := 1 + foreignKeyLen

		whereKeys := make([]string, bindArgsLen)
		bindArgs = make([]interface{}, bindArgsLen)

		i := 0
		{
			pk := schTable.Primary
			f := fieldsMap[schTable.Primary]
			whereKeys[i] = fmt.Sprintf("%s = %s", f.Name, g.RenderBindingValueWithInt(f, bindI))

			bindVal := obj.Get(pk)
			if bindVal == nil {
				return "", nil, &emptyInt, errors.New("dyndao: RenderUpdateWhereClause: missing primary key " + pk)
			}
			bindArgs[i] = bindVal
			i++
			bindI++
		}

		if foreignKeyLen > 0 {
			for _, pk := range schTable.ForeignKeys {
				f := fieldsMap[pk]
				whereKeys[i] = fmt.Sprintf("%s = %s", f.Name, g.RenderBindingValueWithInt(f, bindI))
				bindArgs[i] = obj.Get(pk)
				i++
				bindI++
			}
		}

		whereClause = strings.Join(whereKeys, " AND ")
	}
	return whereClause, bindArgs, &bindI, nil
}

func RenderWhereClause(g *sg.SQLGenerator, schTable *schema.Table, obj *object.Object) (string, []interface{}, error) {
	var whereClause string

	if len(obj.KV) == 0 {
		return "", nil, nil
	}

	whereKeys := make([]string, len(obj.KV))
	bindArgs := make([]interface{}, len(obj.KV))

	i := 0
	bindI := 1
	for k, v := range obj.KV {
		f := schTable.GetColumn(k)
		if f == nil {
			return "", nil, errors.New("dyndao: RenderWhereClause: unknown field " + k + " in table " + obj.Type)
		}
		sqlName := f.Name
		whereKeys[i] = fmt.Sprintf("%s = %s", sqlName, g.RenderBindingValueWithInt(f, bindI))
		switch v.(type) {
		case *object.SQLValue:
			sqlv := v.(*object.SQLValue)
			strV := sqlv.String()
			bindArgs[i] = strV
		default:
			bindArgs[i] = v
		}

		i++
		bindI++
	}
	whereClause = strings.Join(whereKeys, " AND ")
	return whereClause, bindArgs, nil
}

func RenderBindingValueWithInt(f *schema.Column, i int) string {
	return "?"
}
