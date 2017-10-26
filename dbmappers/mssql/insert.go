package mssql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/tidwall/gjson"
)

func RenderInsertValue(f *schema.Column, value interface{}) (interface{}, error) {
	// TODO do we need the schema.Column for more than debugging information?
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
