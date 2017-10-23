package core

import (
	"errors"
	"fmt"
	"os"
	"strings"

	sg "github.com/rbastic/dyndao/sqlgen"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/nils"
)

// TODO: refactor this?
func getRealColumnName(tbl *schema.Table, col string) string {
	if tbl.FieldAliases == nil {
		return col
	}

	realName, ok := tbl.FieldAliases[col]
	if ok {
		return realName
	}
	return col
}

func CoreBindingInsert(g * sg.SQLGenerator, schTable *schema.Table, data map[string]interface{}, identityCol string, fieldsMap map[string]*schema.Field) ([]string, []string, []interface{}) {
	dataLen := len(data)
	bindNames := make([]string, dataLen)
	colNames := make([]string, dataLen)
	bindArgs := make([]interface{}, dataLen)
	i := 0
	for k, v := range data {
		realName := getRealColumnName(schTable, k)
		colNames[i] = realName
		var r string

		if r == "" {
			f, ok := fieldsMap[realName]
			if ok {
				r = g.RenderBindingValue(f)
			} else {
				panic(fmt.Sprintf("coreBindingInsert: Unknown field for key: [%s] realName: [%s] for table %s", k, realName, schTable.Name))
			}
		}

		if v == nil {
			bindNames[i] = r
			bindArgs[i] = v
		} else {
			switch v.(type) {
			case *object.SQLValue:
				sqlv := v.(*object.SQLValue)
				bindNames[i] = sqlv.String()
				bindArgs[i] = nil
			// TODO: dont think this case is necessary...
			case object.SQLValue:
				sqlv := v.(object.SQLValue)
				bindNames[i] = sqlv.String()
				bindArgs[i] = nil
			default:
				bindNames[i] = r
				barg, err := g.RenderInsertValue(fieldsMap[realName], v)
				if err != nil {
					panic(err.Error())
				}
				bindArgs[i] = barg
			}
		}
		i++
	}
	return bindNames, colNames, bindArgs
}

// BindingInsert generates the SQL for a given INSERT statement for oracle with binding parameter values
func BindingInsert(g * sg.SQLGenerator, sch *schema.Schema, table string, data map[string]interface{}) (string, []interface{}, error) {
	if table == "" {
		return "", nil, errors.New("BindingInsert: Empty table name")
	}
	if data == nil {
		return "", nil, errors.New("BindingInsert: Empty data passed")
	}

	schTable := sch.GetTable(table)
	if schTable == nil {
		return "", nil, errors.New("BindingInsert: Table map unavailable for table " + table)
	}

	tableName := schema.GetTableName(schTable.Name, table)

	fieldsMap := schTable.Fields
	if fieldsMap == nil {
		return "", nil, errors.New("BindingInsert: Field map unavailable for table " + table)
	}

	identityCol := schTable.Primary

	bindNames, colNames, bindArgs := g.CoreBindingInsert(g, schTable, data, identityCol, fieldsMap)
	bindArgs = nils.RemoveNilsIfNeeded(bindArgs)

	var sqlStr string

	if schTable.CallerSuppliesPK {
		sqlStr = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			tableName,
			strings.Join(colNames, ","),
			strings.Join(bindNames, ","))
	} else {
		sqlStr = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s /*LASTINSERTID*/ INTO :%s",
			tableName,
			strings.Join(colNames, ","),
			strings.Join(bindNames, ","),
			identityCol,
			identityCol)
	}

	if os.Getenv("DEBUG") != "" {
		fmt.Println("DEBUG: INSERT sqlStr->", sqlStr, "bindArgs->", bindArgs)

	}
	return sqlStr, bindArgs, nil
}
