// Package oraclegen encapsulates an implementation for a given schema attached to
// a generator. This code represents an example implementation for oracle
package oraclegen

import (
	"fmt"
	"strings"

	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

func RenderCreateField(sg * sg.SQLGenerator, f *schema.Field) string {
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

	dataType = mapType(dataType)

	if dataType == "" {
		panic("Empty dataType in renderCreateField for " + f.Name)
	}
	if f.IsIdentity {
		return strings.Join([]string{f.Name, dataType, "GENERATED ALWAYS AS IDENTITY"}, " ")
	}
	return strings.Join([]string{f.Name, dataType, identity, notNull, unique}, " ")
}

func mapType(s string) string {
	// Map 'integer' to 'number' for now for Oracle
	if s == "integer" {
		return "NUMBER"
	}
	if s == "text" {
		return "CLOB"
	}
	return s
}
