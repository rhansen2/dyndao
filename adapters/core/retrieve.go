package core

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

// BindingRetrieve accepts a schema and an object, constructing the appropriate SELECT
// statement to retrieve the object. It will return sqlStr, the EssentialColumns used, and the
// binding where clause.
// DEBUG mode may be turned on by setting an environment parameter, "DEBUG".
func BindingRetrieve(g *sg.SQLGenerator, sch *schema.Schema, obj *object.Object) (string, []string, []interface{}, error) {
	table := obj.Type
	schTable := sch.GetTable(table)
	if schTable == nil {
		return "", nil, nil, errors.New("BindingRetrieve: Table map unavailable for table " + table)
	}

	whereClause, bindWhere, err := g.RenderWhereClause(g, schTable, obj)
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "BindingRetrieve")
	}

	if schTable.EssentialColumns == nil || len(schTable.EssentialColumns) == 0 {
		return "", nil, nil, errors.New("BindingRetrieve: EssentialColumns is empty for table " + table)
	}
	columns := strings.Join(schTable.EssentialColumns, ",")

	whereStr := ""
	if whereClause != "" {
		whereStr = "WHERE"
	}
	tableName := schema.GetTableName(schTable.Name, table)

	sqlStr := fmt.Sprintf("SELECT %s FROM %s %s %s", columns, tableName, whereStr, whereClause)
	return sqlStr, schTable.EssentialColumns, bindWhere, nil
}
