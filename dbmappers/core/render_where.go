package core

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	sg "github.com/rbastic/dyndao/sqlgen"
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

func RenderUpdateWhereClause(g * sg.SQLGenerator, schTable *schema.Table, fieldsMap map[string]*schema.Field, obj *object.Object) (string, []interface{}, error) {
	var bindArgs []interface{}
	var whereClause string

	if len(obj.KV) == 0 {
		return "", nil, nil
	}

	if !schTable.MultiKey {
		f := fieldsMap[schTable.Primary]
		sqlName := f.Name
		whereClause = fmt.Sprintf("%s = %s", sqlName, g.RenderBindingValue(f))
		bindArgs = make([]interface{}, 1)
		bindVal := obj.Get(schTable.Primary)
		if bindVal == nil {
			return "", nil, errors.New("dyndao: RenderUpdateWhereClause: missing primary key " + schTable.Primary)
		}
		bindArgs[0] = bindVal
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
			whereKeys[i] = fmt.Sprintf("%s = %s", f.Name, g.RenderBindingValue(f))

			bindVal := obj.Get(pk)
			if bindVal == nil {
				return "", nil, errors.New("dyndao: RenderUpdateWhereClause: missing primary key " + pk)
			}
			bindArgs[i] = bindVal
			i++
		}

		if foreignKeyLen > 0 {
			for _, pk := range schTable.ForeignKeys {
				f := fieldsMap[pk]
				whereKeys[i] = fmt.Sprintf("%s = %s", f.Name, g.RenderBindingValue(f))
				bindArgs[i] = obj.Get(pk)
				i++
			}
		}

		whereClause = strings.Join(whereKeys, " AND ")
	}
	return whereClause, bindArgs, nil
}

func RenderWhereClause(g * sg.SQLGenerator, schTable *schema.Table, obj *object.Object) (string, []interface{}, error) {
	var whereClause string

	if len(obj.KV) == 0 {
		return "", nil, nil
	}

	whereKeys := make([]string, len(obj.KV))
	bindArgs := make([]interface{}, len(obj.KV))

	i := 0
	for k, v := range obj.KV {
		f := schTable.GetField(k)
		if f == nil {
			return "", nil, errors.New("dyndao: RenderWhereClause: unknown field " + k + " in table " + obj.Type)
		}
		sqlName := f.Name
		whereKeys[i] = fmt.Sprintf("%s = %s", sqlName, g.RenderBindingValue(f))
		bindArgs[i] = v

		i++
	}
	whereClause = strings.Join(whereKeys, " AND ")
	return whereClause, bindArgs, nil
}

func RenderBindingValue(f *schema.Field) string {
	return ":" + f.Name
}

func RenderBindingValueWithInt(f *schema.Field, i int64) string {
	return fmt.Sprintf(":%s%d", f.Name, i)
}
