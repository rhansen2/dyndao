// Package sqlgen encapsulates an implementation for a given schema attached to
// a generator. This code represents an example implementation for sqlite
package sqlgen

import (
	"errors"
	"fmt"

	"github.com/rbastic/dyndao/schema"
)

// CreateTable determines the SQL to create a given table within a schema
func (g Generator) CreateTable(table string) (string, error) {
	tbl, ok := g.Schema.Tables[table]
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

	return "", nil
}

func renderCreateField(f *schema.Field) string {
	dataType := ""
	identity := ""
	nullsAndDefaults := ""

	return fmt.Sprintf("%s %s %s %s", f.Name, dataType, identity, nullsAndDefaults)
}
