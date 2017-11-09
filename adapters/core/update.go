package core

import (
	"github.com/rbastic/dyndao/adapters/common"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

// BindingUpdate generates the SQL for a given UPDATE statement for oracle with binding parameter values
func BindingUpdate(g *sg.SQLGenerator, sch *schema.Schema, obj *object.Object) (string, []interface{}, []interface{}, error) {
	return common.BindingUpdate(g, sch, obj)
}
