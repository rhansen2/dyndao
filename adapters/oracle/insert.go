package oracle

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/tidwall/gjson"
)

func BindingInsertSQL(schTable *schema.Table, tableName string, colNames []string, bindNames []string, identityCol string) string {
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
	return sqlStr
}

func RenderInsertValue(bindI * int, f *schema.Column, value interface{}) (interface{}, error) {
	// TODO do we need the schema.Column for more than debugging information?

	// no need for 'sg', just call the local version
	fName := RenderBindingValueWithIntNoColons(f, *bindI)

	switch value.(type) {
	case string:
		str, ok := value.(string)
		if !ok {
			return "", errors.New("renderInsertValue: unable to turn the value of " + f.Name + " into string")
		}
		return sql.Named(fName, str), nil
	case int32:
		num := value.(int32)
		return sql.Named(fName, string(num)), nil
	case int:
		num := value.(int)
		return sql.Named(fName, num), nil
	case int64:
		num := value.(int64)
		return sql.Named(fName, num), nil
	case uint64:
		num := value.(uint64)
		return sql.Named(fName, fmt.Sprintf("%d", num)), nil
	case float64:
		num := value.(float64)
		if f.IsNumber {
			return sql.Named(fName, int64(num)), nil
		}
		// TODO: when we support more than regular integers, we'll need to care about this more
		return sql.Named(fName, fmt.Sprintf("%f", num)), nil
	case *object.SQLValue:
		val := value.(*object.SQLValue)
		return sql.Named(fName, val.String()), nil
	case object.SQLValue:
		val := value.(object.SQLValue)
		return sql.Named(fName, val.String()), nil
	case gjson.Result:
		panic("gjson.Result is not currently supported for renderInsertValue")
	default:
		return "", fmt.Errorf("renderInsertValue: unknown type %v for the value of (%s, bindName:%s)", reflect.TypeOf(value), f.Name, fName)
	}
}
