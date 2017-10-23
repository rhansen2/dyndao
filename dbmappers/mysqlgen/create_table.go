// Package mysqlgen encapsulates an implementation for a given schema attached to
// a generator. This code represents an example implementation for oracle
package mysqlgen

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rbastic/dyndao/schema"
)

// CreateTable determines the SQL to create a given table within a schema
func (g Generator) CreateTable(s *schema.Schema, table string) (string, error) {
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
		sqlFields[i] = renderCreateField(v)
		i++
	}

	sql := fmt.Sprintf(`CREATE TABLE %s (
	%s
)
`, tableName, strings.Join(sqlFields, ",\n"))

	if os.Getenv("DEBUG") != "" {
		fmt.Println("CreateTable: ", sql)
	}
	return sql, nil
}

func renderCreateField(f *schema.Field) string {
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
