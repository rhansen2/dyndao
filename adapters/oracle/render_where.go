package oracle

import (
	"fmt"


	"github.com/rbastic/dyndao/schema"
)


func RenderBindingValue(f *schema.Column) string {
	return ":"+f.Name
}

func RenderBindingValueWithInt(f *schema.Column, i int64) string {
	return fmt.Sprintf(":%s%d", f.Name, i)
}
