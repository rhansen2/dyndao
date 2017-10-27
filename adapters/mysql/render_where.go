package mysql

import (
	"github.com/rbastic/dyndao/schema"
)

// RenderBindingValue is for binding parameters
func RenderBindingValue(f *schema.Column) string {
	return "?"
}

// RenderBindingValueWithInt is for binding parameters in situations where attaching
// a number as a suffix may be necessary. Not useful for all databases
// (mostly only Oracle, AFAIK at time of writing.)
func RenderBindingValueWithInt(f *schema.Column, i int64) string {
	return "?"
}
