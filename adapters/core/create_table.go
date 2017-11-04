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
func CreateTable(g *sg.SQLGenerator, s *schema.Schema, table string) (string, error) {
	tbl, ok := s.Tables[table]
	if !ok {
		return "", errors.New("dyndao: unknown schema for table with name " + table)
	}
	tableName := schema.GetTableName(tbl.Name, table)
	fieldsMap := tbl.Columns

	sqlColumns := make([]string, len(fieldsMap))
	i := 0
	for _, v := range fieldsMap {
		sqlColumns[i] = g.RenderCreateColumn(g, v)
		i++
	}

	sql := fmt.Sprintf(`CREATE TABLE %s (
	%s
)
`, tableName, strings.Join(sqlColumns, ",\n"))

	if g.Tracing {
		fmt.Printf("dyndao: CreateTable SQL:[%s]\n", sql)
	}

	return sql, nil
}
