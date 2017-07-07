package oraclegen

import (
	"errors"
	"fmt"
	"os"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// BindingDelete generates the appropriate SQL, binding args, and binding where clause parameters
// to execute the requested delete operation. 'obj' is not required to be a
func (g Generator) BindingDelete(sch *schema.Schema, queryVals *object.Object) (string, []interface{}, error) {
	table := queryVals.Type
	schTable, ok := sch.Tables[table]
	if !ok {
		return "", nil, errors.New("BindingDelete: Table map unavailable for table " + table)
	}
	tableName := schema.GetTableName(schTable.Name, table)
	fieldsMap := schTable.Fields
	if fieldsMap == nil {
		return "", nil, errors.New("BindingDelete: Field map unavailable for table " + table)
	}

	whereClause, bindWhere, err := renderWhereClause(schTable, fieldsMap, queryVals)
	if err != nil {
		return "", nil, err
	}

	whereString := "WHERE"
	if len(bindWhere) == 0 {
		whereString = ""
	}
	// TODO: Replicate this fix to sqlite sqlgen
	sqlStr := fmt.Sprintf("DELETE FROM %s %s %s", tableName, whereString, whereClause)
	if os.Getenv("DEBUG") != "" {
		fmt.Println(sqlStr)
	}
	return sqlStr, bindWhere, nil
}
