package oracle

import (
	"fmt"

	"github.com/rbastic/dyndao/schema"
)

func RenderBindingValueWithIntNoColons(f *schema.Column, i int) string {
	return fmt.Sprintf("%s%d", f.Name, i)
}

func RenderBindingValueWithInt(f *schema.Column, i int) string {
	return fmt.Sprintf(":%s%d", f.Name, i)
}
