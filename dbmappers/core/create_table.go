// Package core encapsulates an implementation for a given schema attached to
// a generator. This code represents an example implementation for oracle
package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

// CreateTable determines the SQL to create a given table within a schema
func CreateTable(g * sg.SQLGenerator, s *schema.Schema, table string) (string, error) {
	tbl, ok := s.Tables[table]
	if !ok {
		return "", errors.New("unknown schema for table with name " + table)
	}
	tableName := schema.GetTableName(tbl.Name, table)
	fieldsMap := tbl.Fields

	sqlFields := make([]string, len(fieldsMap))
	i := 0
	// TODO: Have field map in order, or allow one to specify key output order for iterating fields
	// map and generating create SQL....
	for _, v := range fieldsMap {
		sqlFields[i] = g.RenderCreateField(g, v)
		i++
	}

	sql := fmt.Sprintf(`CREATE TABLE %s (
	%s
)
`, tableName, strings.Join(sqlFields, ",\n"))

	return sql, nil
}
