// Package sqlitegen helps with generating SQL statements based on a given schema and additional parameters
package sqlitegen

import (
	"github.com/rbastic/dyndao/schema"
)

// Generator is an empty struct for encapsulating whatever we need for our sql generator ...
type Generator struct {
	Database string
	Name     string
	Schema   *schema.Schema
}

func bindingParam(v string) string {
	return ":" + v
}

// New is our SQLITE code-generator constructor.
func New(db string, name string, sch *schema.Schema) *Generator {
	return &Generator{Database: db, Name: name, Schema: sch}
}

// TODO: Retrieve limited, retrieve paging ... more complexity
