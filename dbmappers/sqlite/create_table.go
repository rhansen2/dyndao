// Package sqlite encapsulates an implementation for a given schema attached to
// a generator. This code represents an example implementation for sqlite
package sqlite

import (
	"fmt"
	"strings"

	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

// TODO: AUTOINCREMENT attribute support.
// See http://www.sqlitetutorial.net/sqlite-autoincrement/

func RenderCreateField(sg *sg.SQLGenerator, f *schema.Field) string {
	dataType := f.DBType
	notNull := ""
	identity := ""
	unique := ""
	if f.IsIdentity {
		identity = "PRIMARY KEY"
	}
	if !f.AllowNull {
		notNull = "NOT NULL"
	}
	if f.IsNumber {
		dataType = f.DBType
	} else {
		if f.Length > 0 {
			dataType = fmt.Sprintf("%s(%d)", f.DBType, f.Length)
		}
	}
	if f.IsUnique {
		unique = "UNIQUE"
	}

	return strings.Join([]string{f.Name, dataType, identity, notNull, unique}, " ")
}
