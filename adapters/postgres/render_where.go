package postgres

import (
	"github.com/rbastic/dyndao/schema"
	"fmt"
)

func RenderBindingValueWithInt(f *schema.Column, i int) string {
	return fmt.Sprintf("$%d", i)
}
