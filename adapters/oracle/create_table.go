// Package oracle encapsulates an implementation for a given schema attached to
// a generator. This code represents an example implementation for oracle
package oracle

import (
	"fmt"
	"strings"

	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

func RenderCreateColumn(sg *sg.SQLGenerator, f *schema.Column) string {
	dataType := f.DBType
	notNull := ""
	identity := ""
	unique := ""
	if f.IsIdentity {
		identity = "PRIMARY KEY"
	}
	if f.AllowNull {
		notNull = "NULL"
	} else {
		notNull = "NOT NULL"
	}

	dataType = mapType(dataType)
	if f.Length > 0 {
		dataType = fmt.Sprintf("%s(%d)", dataType, f.Length)
	}
	if f.IsUnique {
		unique = "UNIQUE"
	}

	if dataType == "" {
		panic("Empty dataType in renderCreateColumn for " + f.Name)
	}
	if f.IsIdentity {
		return strings.Join([]string{f.Name, dataType, "GENERATED ALWAYS AS IDENTITY"}, " ")
	}
	return strings.Join([]string{f.Name, dataType, identity, notNull, unique}, " ")
}

func mapType(s string) string {
	switch s {
	case "integer":
		return "NUMBER"
	case "text":
		return "CLOB"
	case "varchar":
		return "VARCHAR2"
	default:
		return s
	}
}
