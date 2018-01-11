package core

import (
	"errors"
	"fmt"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

// BindingDelete generates the appropriate SQL, binding args, and binding where clause parameters
// to execute the requested delete operation. 'obj' is not required to be a
func BindingDelete(g *sg.SQLGenerator, sch *schema.Schema, queryVals *object.Object) (string, []interface{}, error) {
	table := queryVals.Type
	schTable := sch.GetTable(table)
	if schTable == nil {
		return "", nil, errors.New("BindingDelete: Table map unavailable for table " + table)
	}
	tableName := schema.GetTableName(schTable.Name, table)

	whereClause, bindArgs, err := g.RenderWhereClause(g, schTable, queryVals)
	if err != nil {
		return "", nil, err
	}

	whereString := "WHERE"
	if len(bindArgs) == 0 {
		whereString = ""
	}
	sqlStr := fmt.Sprintf("DELETE FROM %s %s %s", tableName, whereString, whereClause)
	if g.Tracing {
		// TODO: logger interface
		fmt.Println("BindingDelete sqlStr->", sqlStr, "bindArgs->", bindArgs)
	}
	return sqlStr, bindArgs, nil
}
