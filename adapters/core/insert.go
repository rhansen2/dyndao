package core

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	sg "github.com/rbastic/dyndao/sqlgen"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/nils"
	"github.com/tidwall/gjson"
)

// TODO: refactor this?
func getRealColumnName(tbl *schema.Table, col string) string {
	if tbl.ColumnAliases == nil {
		return col
	}

	realName, ok := tbl.ColumnAliases[col]
	if ok {
		return realName
	}
	return col
}

func CoreBindingInsert(g *sg.SQLGenerator, schTable *schema.Table, data map[string]interface{}, identityCol string, fieldsMap map[string]*schema.Column) ([]string, []string, []interface{}) {
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
func BindingInsert(g *sg.SQLGenerator, sch *schema.Schema, table string, data map[string]interface{}) (string, []interface{}, error) {
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

	fieldsMap := schTable.Columns
	if fieldsMap == nil {
		return "", nil, errors.New("BindingInsert: Column map unavailable for table " + table)
	}

	identityCol := schTable.Primary

	bindNames, colNames, bindArgs := g.CoreBindingInsert(g, schTable, data, identityCol, fieldsMap)
	bindArgs = nils.RemoveNilsIfNeeded(bindArgs)

	sqlStr := g.BindingInsertSQL(schTable, tableName, colNames, bindNames, identityCol)

	return sqlStr, bindArgs, nil
}

func BindingInsertSQL(schTable *schema.Table, tableName string, colNames []string, bindNames []string, identityCol string) string {
	var sqlStr string
	sqlStr = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(colNames, ","),
		strings.Join(bindNames, ","))
	return sqlStr
}

func RenderInsertValue(f *schema.Column, value interface{}) (interface{}, error) {
	// TODO do we need the schema.Column for more than debugging information?
	switch value.(type) {
	case string:
		str, ok := value.(string)
		if !ok {
			return "", errors.New("renderInsertValue: unable to turn the value of " + f.Name + " into string")
		}
		return str, nil
	case int32:
		num := value.(int32)
		return string(num), nil
	case int:
		num := value.(int)
		return num, nil
	case int64:
		num := value.(int64)
		return num, nil
	case uint64:
		num := value.(uint64)
		return fmt.Sprintf("%d", num), nil
	case float64:
		num := value.(float64)
		if f.IsNumber {
			return int64(num), nil
		}
		// TODO: when we support more than regular integers, we'll need to care about this more
		return fmt.Sprintf("%f", num), nil
	case *object.SQLValue:
		val := value.(*object.SQLValue)
		return val.String(), nil
	case object.SQLValue:
		val := value.(object.SQLValue)
		return val.String(), nil
	case gjson.Result:
		panic("gjson.Result is not currently supported for renderInsertValue")
	default:
		return "", fmt.Errorf("renderInsertValue: unknown type %v for the value of %s", reflect.TypeOf(value), f.Name)
	}
}
