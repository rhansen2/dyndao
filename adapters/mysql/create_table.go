// Package mysql encapsulates an implementation for a given schema attached to
// a generator. This code represents an example implementation for oracle
package mysql

import (
	"fmt"
	"strings"

	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

func RenderCreateColumn(sg *sg.SQLGenerator, f *schema.Column) string {
	dataType := f.DBType
	var notNull string
	identity := ""
	unique := ""

	if f.IsIdentity {
		identity = "PRIMARY KEY AUTO_INCREMENT"
		f.AllowNull = false
	}

	if !f.AllowNull {
		notNull = "NOT NULL"
	} else {
		notNull = "NULL"
	}

	if f.Length > 0 {
		dataType = fmt.Sprintf("%s(%d)", f.DBType, f.Length)
	}

	if f.IsUnique {
		unique = "UNIQUE"
	}

	dataType = mapType(dataType)

	if dataType == "" {
		panic("Empty dataType in renderCreateColumn for " + f.Name)
	}
	return strings.Join([]string{f.Name, dataType, identity, notNull, unique}, " ")
}

func mapType(s string) string {
	// Map 'integer' to 'int(11)' for now for MySQL
	if s == "integer" {
		return "int(11)"
	}
	// no need to map text type
	return s
}
