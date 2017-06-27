// Package sqlitegen helps with generating SQL statements based on a given schema and additional parameters
package sqlitegen

import "github.com/rbastic/dyndao/schema"

// Generator is an empty struct for encapsulating whatever we need for our sql generator ...
type Generator struct {
	Database string
	Name     string
}

func bindingParam(v string) string {
	return ":" + v
}

// New is our SQLITE code-generator constructor.
func New(name string, sch *schema.Schema) *Generator {
	return &Generator{Database: "sqlite", Name: name}
}
