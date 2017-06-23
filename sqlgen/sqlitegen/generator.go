// Package sqlitegen helps with generating SQL statements based on a given schema and additional parameters
package sqlitegen

import (
	"fmt"
	"github.com/rbastic/dyndao/schema"
	"strings"
)

// Generator is an empty struct for encapsulating whatever we need for our sql generator ...
type Generator struct {
	Database string
	Name     string
	Schema   *schema.Schema
}

func bindingParam(v string) string {
	return ":" + v
}

func New( db string, name string, sch * schema.Schema ) *Generator {
	return &Generator{ Database: db, Name: name, Schema: sch }
}

// Retrieve is the generic retrieve function to retrieve a single or multiple records
// It could be combined with a prepared statement for speed, or used individually (on an
// ad-hoc basis)
func (g Generator) Retrieve(table string, fields []string, pkValues map[string]string) string {
	fieldsStr := strings.Join(fields, ",")

	valuesAry := make([]string, len(pkValues))
	i := 0
	for k, v := range pkValues {
		valuesAry[i] = fmt.Sprintf("%s = %s", k, bindingParam(v))
		i++
	}
	pkValuesStr := strings.Join(valuesAry, ",")

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s", table, fieldsStr, pkValuesStr)

	return sql
}

// TODO: Retrieve limited, retrieve paging ... more complexity
