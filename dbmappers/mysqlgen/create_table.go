// Package mysqlgen encapsulates an implementation for a given schema attached to
// a generator. This code represents an example implementation for oracle
package mysqlgen

import (
	"fmt"
	"strings"

	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

func RenderCreateField(sg * sg.SQLGenerator, f *schema.Field) string {
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

	if dataType == "" {
		panic("Empty dataType in renderCreateField for " + f.Name)
	}
	return strings.Join([]string{f.Name, dataType, identity, notNull, unique}, " ")
}
