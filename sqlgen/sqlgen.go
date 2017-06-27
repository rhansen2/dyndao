package sqlgen

import (
	"errors"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/dyndao/sqlgen/sqlitegen"
)

type GeneratorImp interface {
	BindingInsert(table string, data map[string]interface{}) (string, []interface{}, error)
	Insert(table string, data map[string]interface{}) (string, error)
	BindingUpdate(sch *schema.Schema, obj *object.Object) (string, []interface{}, []interface{}, error)
	BindingRetrieve(sch *schema.Schema, obj *object.Object) (string, []interface{}, error)
}

// Generator is an empty struct for encapsulating whatever we need for our sql generator ...
type Generator struct {
	Database string
	Name     string
	Schema   *schema.Schema
}

// New is our generic sql generator constructor
func New(db string, name string, sch *schema.Schema) (interface{}, error) {
	switch db {
	case "sqlite":
		//	fallthrough
		//case "oracle":
		// TODO: fix testName as a parameter
		return sqlitegen.New(db, "testName", sch), nil
	default:
		return nil, errors.New("sqlgen: Unrecognized database type " + db)
	}
}
