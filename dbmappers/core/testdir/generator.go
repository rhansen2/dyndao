// Package oraclegen helps with generating SQL statements based on a given schema and additional parameters
package core

import "github.com/rbastic/dyndao/schema"

// Generator is an empty struct for encapsulating whatever we need for our sql generator ...
type Generator struct {
	Database string
	Name     string
}

// New is our Oracle 'code-generator' constructor.
func New(name string, sch *schema.Schema) *Generator {
	return &Generator{Database: "oracle", Name: name}
}

// FixLastInsertIDbug is a nasty hack to deal with some bugs I found in rana's
// ora.v4 oracle driver. FIXME Add more information here.
func (g Generator) FixLastInsertIDbug() bool {
	return true
}