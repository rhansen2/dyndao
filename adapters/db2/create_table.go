package db2

import (
	"fmt"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
	"strings"
)

func RenderCreateColumn(sg *sg.SQLGenerator, f *schema.Column) string {
	// TODO: Code a real TypeMapper
	f.DBType = mapType(f.DBType)
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
		return strings.Join([]string{f.Name, dataType, "NOT NULL GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1)"}, " ")
	}
	return strings.Join([]string{f.Name, dataType, identity, notNull, unique}, " ")
}

func mapType(s string) string {
	switch s {
	case "INTEGER":
		return "INT"
	case "text":
		return "CLOB"
	case "TEXT":
		return "CLOB"
	case "blob":
		return "CLOB"
	case "BLOB":
		return "CLOB"
	default:
		return s
	}
}
