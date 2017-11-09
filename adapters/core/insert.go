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

func bindingValueHelper( g *sg.SQLGenerator, fieldsMap map[string]*schema.Column, realName string, bindI *int, k string, schTable * schema.Table) string {
	f, ok := fieldsMap[realName]
	if !ok {
		panic(fmt.Sprintf("coreBindingInsert: Unknown field for key: [%s] realName: [%s] for table %s", k, realName, schTable.Name))
	}
	r := g.RenderBindingValueWithInt(f, *bindI)
	return r
}

func CoreBindingInsert(g *sg.SQLGenerator, schTable *schema.Table, data map[string]interface{}, identityCol string, fieldsMap map[string]*schema.Column) ([]string, []string, []interface{}) {
	dataLen := len(data)
	bindNames := make([]string, dataLen)
	colNames := make([]string, dataLen)
	bindArgs := make([]interface{}, dataLen)
	i := 0
	bindI := 1
	for k, v := range data {
		realName := schTable.GetColumnName(k)
		colNames[i] = realName

		switch v.(type) {
		case *object.SQLValue:
			sqlv := v.(*object.SQLValue)
			bindNames[i] = sqlv.String()
			bindArgs[i] = nil
/*		// TODO: dont think this case is necessary...
		case object.SQLValue:
			sqlv := v.(object.SQLValue)
			bindNames[i] = sqlv.String()
			bindArgs[i] = nil
*/
		case nil:
			bindNames[i] = bindingValueHelper(g, fieldsMap, realName, &bindI, k, schTable)
			bindI++
			bindArgs[i] = v
		default:
			bindNames[i] = bindingValueHelper(g, fieldsMap, realName, &bindI, k, schTable)
			barg, err := g.RenderInsertValue(&bindI, fieldsMap[realName], v)
			if err != nil {
				panic(err)
			}
			bindI++
			bindArgs[i] = barg
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

func RenderInsertValue(bindI * int, f *schema.Column, value interface{}) (interface{}, error) {
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
