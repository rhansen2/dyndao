package oraclegen

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/nils"
	"github.com/tidwall/gjson"
)

func (g Generator) setPKifNeeded(data map[string]interface{}, identityCol string) {
	if g.CallerSuppliesPK {
		// NOTE: Currently we support this as a SYS_GUID implementation,
		// and technically, /we/ aren't the caller ... maybe I will rename
		// this as ORMSuppliesPK ... but I don't have another reference
		// point to compare against yet. YAGNI applies.
		_, ok := data[identityCol]
		if !ok {
			data[identityCol] = object.NewSQLValue("SYS_GUID()")
		}
	}
}

func (g Generator) coreBindingInsert(schTable *schema.Table, data map[string]interface{}, identityCol string, fieldsMap map[string]*schema.Field) ([]string, []string, []interface{}) {
	dataLen := len(data)
	bindNames := make([]string, dataLen)
	colNames := make([]string, dataLen)
	bindArgs := make([]interface{}, dataLen)
	i := 0
	for k, v := range data {
		realName := getRealColumnName(schTable, k)
		colNames[i] = realName
		var r string

		if g.CallerSuppliesPK && k == identityCol {
			switch typ := v.(type) {
			case int64:
				r = string(v.(int64))
			case string:
				r = v.(string)
			case *object.SQLValue:
				thing := v.(*object.SQLValue)
				r = thing.Value
				v = nil
			default:
				panic(fmt.Sprintf("coreBindingInsert: Unknown type [%v] in switch", typ))
			}
		}
		if r == "" {
			f, ok := fieldsMap[realName]
			if ok {
				r = g.RenderBindingValue(f)
			} else {
				panic(fmt.Sprintf("coreBindingInsert: Unknown field for key: [%s] realName: [%s] for table %s", k, realName, schTable.Name))
			}
		}
		bindNames[i] = r
		if v == nil {
			bindArgs[i] = v
		} else {
			barg, err := renderInsertValue(fieldsMap[realName], v)
			if err != nil {
				panic(err.Error())
			}
			bindArgs[i] = barg

		}
		i++
	}
	return bindNames, colNames, bindArgs
}

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

// BindingInsert generates the SQL for a given INSERT statement for oracle with binding parameter values
func (g Generator) BindingInsert(sch *schema.Schema, table string, data map[string]interface{}) (string, []interface{}, error) {
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
	g.setPKifNeeded(data, identityCol)

	bindNames, colNames, bindArgs := g.coreBindingInsert(schTable, data, identityCol, fieldsMap)
	bindArgs = nils.RemoveNilsIfNeeded(bindArgs)

	var sqlStr string
	if g.CallerSuppliesPK {
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

func renderInsertValue(f *schema.Field, value interface{}) (interface{}, error) {
	// TODO do we need the schema.Field for more than debugging information?
	switch value.(type) {
	case string:
		str, ok := value.(string)
		if !ok {
			return "", errors.New("renderInsertValue: unable to turn the value of " + f.Name + " into string")
		}
		return sql.Named(f.Name, str), nil
	case int32:
		num := value.(int32)
		return sql.Named(f.Name, string(num)), nil
	case int:
		num := value.(int)
		return sql.Named(f.Name, num), nil
	case int64:
		num := value.(int64)
		return sql.Named(f.Name, num), nil
	case uint64:
		num := value.(uint64)
		return sql.Named(f.Name, fmt.Sprintf("%d", num)), nil
	case float64:
		num := value.(float64)
		if f.IsNumber {
			return sql.Named(f.Name, int64(num)), nil
		}
		// TODO: when we support more than regular integers, we'll need to care about this more
		return sql.Named(f.Name, fmt.Sprintf("%f", num)), nil
	case *object.SQLValue:
		val := value.(*object.SQLValue)
		return sql.Named(f.Name, val.String()), nil
	case object.SQLValue:
		val := value.(object.SQLValue)
		return sql.Named(f.Name, val.String()), nil
	case gjson.Result:
		panic("gjson.Result is not currently supported for renderInsertValue")
	default:
		return "", fmt.Errorf("renderInsertValue: unknown type %v for the value of %s", reflect.TypeOf(value), f.Name)
	}
}
