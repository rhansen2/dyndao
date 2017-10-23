package mysqlgen

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// BindingRetrieve accepts a schema and an object, constructing the appropriate SELECT
// statement to retrieve the object. It will return sqlStr, the EssentialFields used, and the
// binding where clause.
// DEBUG mode may be turned on by setting an environment parameter, "DEBUG".
// TODO: We may consider using a different name in the future.
func (g Generator) BindingRetrieve(sch *schema.Schema, obj *object.Object) (string, []string, []interface{}, error) {
	table := obj.Type // TODO: we may want to map this
	schTable := sch.GetTable(table)
	if schTable == nil {
		return "", nil, nil, errors.New("BindingRetrieve: Table map unavailable for table " + table)
	}

	whereClause, bindWhere, err := g.renderWhereClause(schTable, obj)
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "BindingRetrieve")
	}

	if schTable.EssentialFields == nil || len(schTable.EssentialFields) == 0 {
		return "", nil, nil, errors.New("BindingRetrieve: EssentialFields is empty for table " + table)
	}
	columns := strings.Join(schTable.EssentialFields, ",")

	whereStr := ""
	if whereClause != "" {
		whereStr = "WHERE"
	}
	tableName := schema.GetTableName(schTable.Name, table)

	sqlStr := fmt.Sprintf("SELECT %s FROM %s %s %s", columns, tableName, whereStr, whereClause)
	return sqlStr, schTable.EssentialFields, bindWhere, nil
}
