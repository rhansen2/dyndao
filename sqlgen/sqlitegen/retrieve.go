package sqlitegen

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// BindingRetrieve is a simple binding retrieve.
func (g Generator) BindingRetrieve(sch *schema.Schema, obj *object.Object) (string, []interface{}, error) {
	table := obj.Type // TODO: we may want to map this
	schTable, ok := g.Schema.Tables[table]
	if !ok {
		return "", nil, errors.New("BindingRetrieve: Table map unavailable for table " + table)
	}

	fieldsMap := schTable.Fields
	if fieldsMap == nil {
		return "", nil, errors.New("BindingRetrieve: Field map unavailable for table " + table)
	}

	whereClause, bindWhere, err := renderWhereClause(schTable, fieldsMap, obj)
	if err != nil {
		return "", nil, errors.Wrap(err, "BindingRetrieve: renderWhereClause")
	}

	if schTable.EssentialFields == nil || len(schTable.EssentialFields) == 0 {
		return "", nil, errors.New("BindingRetrieve: EssentialFields is empty")
	}
	columns := strings.Join(schTable.EssentialFields, ",")

	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s", columns, table, whereClause)
	//fmt.Println(sqlStr)

	return sqlStr, bindWhere, nil
}

func renderWhereClause(schTable *schema.Table, fieldsMap map[string]*schema.Field, obj *object.Object) (string, []interface{}, error) {
	whereClause := ""
	var bindArgs []interface{}

	if !schTable.MultiKey {
		f := fieldsMap[schTable.Primary]
		sqlName := f.Name
		whereClause = fmt.Sprintf("%s = %s", sqlName, renderBindingUpdateValue(f))
		bindArgs = make([]interface{}, 1)
		bindArgs[0] = obj.Get(schTable.Primary)
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
			bindArgs[i] = obj.Get(pk)
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
