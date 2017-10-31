// Package mysql encapsulates an implementation for a given schema attached to
// a generator. This code represents an example implementation for oracle
package mysql

import (
	"github.com/rbastic/dyndao/adapters/common"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

func RenderCreateColumn(sg *sg.SQLGenerator, f *schema.Column) string {
	if f.IsIdentity {
		f.AllowNull = false
	}

	return common.RenderCreateColumn(sg, f, "PRIMARY KEY AUTO_INCREMENT", mapType)
}

func mapType(s string) string {
	// Map 'integer' to 'int(11)' for now for MySQL
	if s == "INTEGER" {
		return "INT(11)"
	}
	// no need to map text type
	return s
}
