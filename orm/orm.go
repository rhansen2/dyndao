// Package orm is only a test ORM package meant to demonstrate how to code your
// own dynamic ORM layer.
package orm

import (
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/dyndao/sqlitegen"
)

// Insert returns an integer of rows affected and possibly an error.
func Insert(sch * schema.Schema, obj * object.Object, sql string) (int, error) {
	return 0, nil
}

// Update returns an integer of rows affected and possibly an error.
func Update(sch * schema.Schema, obj * object.Object, sql string) (int, error) {
	return 0, nil
}

// Retrieve returns a fully fleshened object, given a relevant schema, table,
// and primary key fields. It may return an error.
func Retrieve( sch * schema.Schema, tbl string, pkFields []string) (*object.Object, error) {

}
