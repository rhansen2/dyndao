// Package oraclegen helps with generating SQL statements based on a given schema and additional parameters
package oraclegen

import "github.com/rbastic/dyndao/schema"

// Generator is an empty struct for encapsulating whatever we need for our sql generator ...
type Generator struct {
	Database         string
	Name             string
	CallerSuppliesPK bool
}

func bindingParam(v string) string {
	return ":" + v
}

// New is our Oracle 'code-generator' constructor.
func New(name string, sch *schema.Schema, callerSuppliesPK bool) *Generator {
	return &Generator{Database: "oracle", Name: name, CallerSuppliesPK: callerSuppliesPK}
}

func (g Generator) FixLastInsertIDbug() bool {
	return true
}
