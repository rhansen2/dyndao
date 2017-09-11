package sqlitegen

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
	table := obj.Type
	if table == "" {
		return "", nil, nil, errors.New("BindingRetrieve: Empty object type received")
	}

	schTable := sch.GetTable(table)
	if schTable == nil {
		return "", nil, nil, errors.New("BindingRetrieve: Table map unavailable for table " + table)
	}

	whereClause, bindWhere, err := renderWhereClause(schTable, obj)
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "BindingRetrieve")
	}

	columns := strings.Join(schTable.EssentialFields, ",")

	whereStr := ""
	if whereClause != "" {
		whereStr = "WHERE"
	}
	sqlStr := fmt.Sprintf("SELECT %s FROM %s %s %s", columns, schTable.Name, whereStr, whereClause)
	return sqlStr, schTable.EssentialFields, bindWhere, nil
}
