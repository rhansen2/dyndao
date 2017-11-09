package postgres

import (
	"github.com/rbastic/dyndao/adapters/common"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

func BindingUpdate(g *sg.SQLGenerator, sch *schema.Schema, obj *object.Object) (string, []interface{}, []interface{}, error) {

	sqlStr, bindArgs, bindWhere, err := common.BindingUpdate(g, sch, obj)
	return sqlStr, bindWhere, bindArgs, err
}
