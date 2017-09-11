// Package sqlitegen encapsulates an implementation for a given schema attached to
// a generator. This code represents an example implementation for sqlite
package sqlitegen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rbastic/dyndao/schema"
)

// CreateTable determines the SQL to create a given table within a schema
func (g Generator) CreateTable(s *schema.Schema, table string) (string, error) {
	tbl := s.GetTable(table)
	if tbl == nil {
		return "", errors.New("unknown schema for table with name " + table)
	}
	tableName := schema.GetTableName(tbl.Name, table)
	fieldsMap := tbl.Fields

	sqlFields := make([]string, len(fieldsMap))
	i := 0
	// TODO: Have field map in order, or allow one to specify key output order for iterating fields
	// map and generating create SQL....
	for _, v := range fieldsMap {
		sqlFields[i] = renderCreateField(v)
		i++
	}

	sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	%s
)
`, tableName, strings.Join(sqlFields, ",\n"))

	return sql, nil
}

func renderCreateField(f *schema.Field) string {
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
