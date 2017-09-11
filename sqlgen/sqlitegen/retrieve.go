package sqlitegen

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// BindingRetrieve accepts a schema and an object, constructing the appropriate SELECT
// statement to retrieve the object. It will return sqlStr, the EssentialFields used, and the
// binding where clause.
// DEBUG mode may be turned on by setting an environment parameter, "DEBUG".
// TODO: We may consider using a different name in the future.
func (g Generator) BindingRetrieve(sch *schema.Schema, obj *object.Object) (string, []string, []interface{}, error) {
	table := obj.Type // TODO: we may want to map this
	schTable := sch.GetTable(table)
	if schTable == nil {
		return "", nil, nil, errors.New("BindingRetrieve: Table map unavailable for table " + table)
	}

	whereClause, bindWhere, err := renderWhereClause(schTable, obj)
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "BindingRetrieve")
	}

	if schTable.EssentialFields == nil || len(schTable.EssentialFields) == 0 {
		return "", nil, nil, errors.New("BindingRetrieve: EssentialFields is empty for table " + table)
	}
	columns := strings.Join(schTable.EssentialFields, ",")

	whereStr := ""
	if whereClause != "" {
		whereStr = "WHERE"
	}
	tableName := schema.GetTableName(schTable.Name, table)

	sqlStr := fmt.Sprintf("SELECT %s FROM %s %s %s", columns, tableName, whereStr, whereClause)
	return sqlStr, schTable.EssentialFields, bindWhere, nil
}

func renderUpdateWhereClause(schTable *schema.Table, fieldsMap map[string]*schema.Field, obj *object.Object) (string, []interface{}, error) {
	var bindArgs []interface{}
	var whereClause string

	if len(obj.KV) == 0 {
		return "", nil, nil
	}

	if !schTable.MultiKey {
		f := fieldsMap[schTable.Primary]
		sqlName := f.Name
		whereClause = fmt.Sprintf("%s = %s", sqlName, renderBindingUpdateValue(f))
		bindArgs = make([]interface{}, 1)
		bindVal := obj.Get(schTable.Primary)
		if bindVal == nil {
			return "", nil, errors.New("renderWhereClause: missing primary key " + schTable.Primary)
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
			whereKeys[i] = fmt.Sprintf("%s = %s", f.Name, renderBindingUpdateValue(f))

			bindVal := obj.Get(pk)
			if bindVal == nil {
				return "", nil, errors.New("renderWhereClause: missing primary key " + pk)
			}
			bindArgs[i] = bindVal
			i++
		}

		if foreignKeyLen > 0 {
			for _, pk := range schTable.ForeignKeys {
				f := fieldsMap[pk]
				whereKeys[i] = fmt.Sprintf("%s = %s", f.Name, renderBindingUpdateValue(f))
				bindArgs[i] = obj.Get(pk)
				i++
			}
		}

		whereClause = strings.Join(whereKeys, " AND ")
	}
	return whereClause, bindArgs, nil
}

func renderWhereClause(schTable *schema.Table, obj *object.Object) (string, []interface{}, error) {
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
			return "", nil, errors.New("renderWhereClause: unknown field " + k + " in table " + obj.Type)
		}
		sqlName := f.Name
		whereKeys[i] = fmt.Sprintf("%s = %s", sqlName, renderBindingUpdateValue(f))
		bindArgs[i] = v

		i++
	}
	whereClause = strings.Join(whereKeys, " AND ")
	return whereClause, bindArgs, nil
}
