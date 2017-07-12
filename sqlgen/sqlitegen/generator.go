// Package sqlitegen helps with generating SQL statements based on a given schema and additional parameters
package sqlitegen

import "github.com/rbastic/dyndao/schema"

// Generator is an empty struct for encapsulating whatever we need for our sql generator ...
type Generator struct {
	Database         string
	Name             string
	CallerSuppliesPK bool
}

// New is our SQLITE code-generator constructor.
func New(name string, sch *schema.Schema, callerSuppliesPK bool) *Generator {
	return &Generator{Database: "oracle", Name: name, CallerSuppliesPK: callerSuppliesPK}
}

// FixLastInsertIDbug is a nasty hack to deal with some bugs I found in rana's
// ora.v4 oracle driver. FIXME Add more information here.
func (g Generator) FixLastInsertIDbug() bool {
	return false
}

// CallerSuppliesPrimaryKey is a boolean accessor that returns the current global
// primary key mode. FIXME: Add more notes on this.
func (g Generator) CallerSuppliesPrimaryKey() bool {
	return g.CallerSuppliesPK
}
