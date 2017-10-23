package sqlitegen

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/rbastic/dyndao/schema"
)

func quotedString(value string) string {
	return fmt.Sprintf(`"%s"`, value)
}

func RenderInsertValue(f *schema.Field, value interface{}) (interface{}, error) {
	// TODO do we need the schema.Field for more than debugging information?
	switch v := value.(type) {
	case string:
		if v == "" {
			return "", errors.New("dyndao: RenderInsertValue: unable to turn the value of " + f.Name + " into string")
		}
		return quotedString(v), nil
	case int32:
		num := value.(int32)
		return string(num), nil
	case int:
		num := value.(int)
		return strconv.Itoa(num), nil
	case int64:
		num := value.(int64)
		return string(num), nil
	default:
		return "", errors.New("dyndao: RenderInsertValue: unknown type " + fmt.Sprintf("%v", v) + " for the value of " + f.Name)

	}
}

func RenderBindingValue(f *schema.Field) string {
	return ":" + f.Name
}

func RenderBindingValueWithInt(f *schema.Field, i int64) string {
	return fmt.Sprintf(":%s%d", f.Name, i)
}
