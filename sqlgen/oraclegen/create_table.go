// Package oraclegen encapsulates an implementation for a given schema attached to
// a generator. This code represents an example implementation for sqlite
package oraclegen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rbastic/dyndao/schema"
)

// CreateTable determines the SQL to create a given table within a schema
func (g Generator) CreateTable(s *schema.Schema, table string) (string, error) {
	tbl, ok := s.Tables[table]
	if !ok {
		return "", errors.New("unknown schema for table with name " + table)
	}

	fieldsMap := tbl.Fields

	sqlFields := make([]string, len(fieldsMap))
	i := 0
	// TODO: Have field map in order, or allow one to specify key output order for iterating fields
	// map and generating create SQL....
	for _, v := range fieldsMap {
		sqlFields[i] = renderCreateField(v)
		i++
	}

	sql := fmt.Sprintf(`CREATE TABLE %s (
	%s
)
`, table, strings.Join(sqlFields, ",\n"))

	return sql, nil
}

func renderCreateField(f *schema.Field) string {
	dataType := ""
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
		dataType = f.DBType // not relevant here?
	} else {
		if f.Length > 0 {
			// TODO: Needed?
		}
		dataType = f.DBType
	}
	if f.IsUnique {
		unique = "UNIQUE"
	}

	// Map 'integer' to 'number' for now for Oracle
	if dataType == "integer" {
		dataType = "NUMBER"
	}
	if dataType == "text" {
		dataType = "CLOB"
	}
	if f.IsIdentity {
		return strings.Join([]string{f.Name, dataType, "GENERATED ALWAYS AS IDENTITY"}, " ")
	}
	return strings.Join([]string{f.Name, dataType, identity, notNull, unique}, " ")
}
