package postgres

import (
	"fmt"
	"github.com/rbastic/dyndao/schema"
	"strings"
	sg "github.com/rbastic/dyndao/sqlgen"
)

var (
	identityStr = "SERIAL PRIMARY KEY"
)

func RenderCreateColumn(sg *sg.SQLGenerator, f *schema.Column) string {
	dataType := strings.ToUpper(f.DBType)

	notNull := ""
	identity := ""
	unique := ""

	if f.IsIdentity {
		identity = identityStr
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

	if f.IsIdentity {
		return strings.Join([]string{f.Name, identity, notNull, unique}, " ")
	}
	return strings.Join([]string{f.Name, dataType, identity, notNull, unique}, " ")
}

func mapType(s string) string {
	// Map 'integer' to 'number' for now for Oracle
	if s == "INTEGER" {
		return "INT"
	}
	if s == "TEXT" {
		return "TEXT"
	}
	// HRM
	if s == "BLOB" || s == "CLOB" {
		return "TEXT"
	}
	return s
}
