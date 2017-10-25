package sqlite

import (
	"errors"
	"fmt"

	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/dyndao/object"
	"github.com/tidwall/gjson"
)

func quotedString(value string) string {
	return fmt.Sprintf(`"%s"`, value)
}

func RenderInsertValue(f *schema.Column, value interface{}) (interface{}, error) {
	// TODO do we need the schema.Column for more than debugging information?
	switch v := value.(type) {
	case string:
		if v == "" {
			return "", errors.New("dyndao: RenderInsertValue: unable to turn the value of " + f.Name + " into string")
		}
		return v, nil
	case int32:
		num := value.(int32)
		return num, nil
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
		return "", errors.New("dyndao: RenderInsertValue: unknown type " + fmt.Sprintf("%v", v) + " for the value of " + f.Name)

	}
}

func RenderBindingValue(f *schema.Column) string {
	return ":" + f.Name
}

func RenderBindingValueWithInt(f *schema.Column, i int64) string {
	return fmt.Sprintf(":%s%d", f.Name, i)
}
