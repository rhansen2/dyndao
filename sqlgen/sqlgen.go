package sqlgen

import (
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
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
