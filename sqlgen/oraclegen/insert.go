package oraclegen

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rbastic/dyndao/schema"
)

// TODO: Refactor common code... later. A lot of overall work remains.

// BindingInsert generates the SQL for a given INSERT statement for oracle with binding parameter values
func (g Generator) BindingInsert(sch *schema.Schema, table string, data map[string]interface{}) (string, []interface{}, error) {
	if table == "" {
		return "", nil, errors.New("BindingInsert: Empty table name")
	}
	if data == nil {
		return "", nil, errors.New("BindingInsert: Empty data passed")
	}

	schTable, ok := sch.Tables[table]
	if !ok {
		return "", nil, errors.New("BindingInsert: Table map unavailable for table " + table)
	}
	fieldsMap := schTable.Fields
	if fieldsMap == nil {
		return "", nil, errors.New("BindingInsert: Field map unavailable for table " + table)
	}
	identityCol := schTable.Primary

	dataLen := len(data)
	bindNames := make([]string, dataLen)
	colNames := make([]string, dataLen)
	bindArgs := make([]interface{}, dataLen)
	i := 0

	for k, v := range data {
		colNames[i] = fmt.Sprintf(`%s`, k)
		//fmt.Println("k=", k, "fieldsMap[k]=", fieldsMap[k], "v=", v)
		r := renderBindingInsertValue(fieldsMap[k])
		bindNames[i] = fmt.Sprintf(`%s`, r)
		bindArgs[i] = v
		i++
	}
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s /*LASTINSERTID*/ INTO :%s", table, strings.Join(colNames, ","), strings.Join(bindNames, ","), identityCol, identityCol)
	return sqlStr, bindArgs, nil
}

func quotedString(value string) string {
	return fmt.Sprintf(`"%s"`, string(value))
}

func renderBindingInsertValue(f *schema.Field) string {
	return ":" + f.Name
}

func renderBindingRetrieve(f *schema.Field) string {
	return renderBindingUpdateValue(f)
}

func renderInsertValue(f *schema.Field, value interface{}) (string, error) {
	// TODO do we need the schema.Field for more than debugging information?
	switch v := value.(type) {
	case string:
		str := string(v)
		if str == "" {
			return "", errors.New("renderInsertField: unable to turn the value of " + f.Name + " into string")
		}
		return quotedString(str), nil
	case int32:
		num := value.(int32)
		return string(num), nil
	case int:
		num := value.(int)
		return strconv.Itoa(num), nil
	case int64:
		num := value.(int64)
		return string(num), nil
	default:
		return "", errors.New("renderInsertField: unknown type " + fmt.Sprintf("%v", v) + " for the value of " + f.Name)

	}
	//return "", nil
}

// TODO: InsertBinding
