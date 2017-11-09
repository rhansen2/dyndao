// Package sqlite encapsulates an implementation for a given schema attached to
// a generator. This code represents an example implementation for sqlite
package sqlite

import (
	"github.com/rbastic/dyndao/adapters/common"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

// TODO: AUTOINCREMENT attribute support.
// See http://www.sqlitetutorial.net/sqlite-autoincrement/

func RenderCreateColumn(sg *sg.SQLGenerator, f *schema.Column) string {
	return common.RenderCreateColumn(sg, f, "PRIMARY KEY", mapType)
}

func mapType(s string) string {
	switch s {
	case "TIMESTAMP":
		return "DATETIME"
	default:
		return s
	}
}
