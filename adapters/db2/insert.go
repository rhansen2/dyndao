package db2

import (
	"fmt"
	"strings"

	"github.com/rbastic/dyndao/schema"
)

func BindingInsertSQL(schTable *schema.Table, tableName string, colNames []string, bindNames []string, identityCol string) string {
	var sqlStr string
	//if schTable.CallerSuppliesPK {
		sqlStr = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			tableName,
			strings.Join(colNames, ","),
			strings.Join(bindNames, ","))
	/*} else {
		sqlStr = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING INTO %s",
			tableName,
			strings.Join(colNames, ","),
			strings.Join(bindNames, ","),
			identityCol)
	}*/
	return sqlStr
}

