package postgres

import (
	"fmt"
	"github.com/rbastic/dyndao/schema"
)

func RenderBindingValueWithInt(f *schema.Column, i int) string {
	return fmt.Sprintf("$%d", i)
}
