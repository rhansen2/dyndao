package sqlitegen

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rbastic/dyndao/schema"
)

// TODO: Refactor common code... later. A lot of overall work remains.

// BindingInsert generates the SQL for a given INSERT statement for SQLite with binding parameter values
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
	dataLen := len(data)
	//fmt.Println(dataLen)
	bindNames := make([]string, dataLen)
	colNames := make([]string, dataLen)
	bindArgs := make([]interface{}, dataLen)
	i := 0

	for k, v := range data {
		colNames[i] = k
		//fmt.Println("k=", k, "fieldsMap[k]=", fieldsMap[k], "v=", v)
		r := renderBindingInsertValue(fieldsMap[k])
		bindNames[i] = r
		bindArgs[i] = v
		i++
	}
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(colNames, ","), strings.Join(bindNames, ","))
	return sqlStr, bindArgs, nil

}

// Insert generates the SQL for a given INSERT statement for SQLite
func (g Generator) Insert(sch *schema.Schema, table string, data map[string]interface{}) (string, error) {
	if table == "" {
		return "", errors.New("Insert: empty table name")
	}
	if data == nil {
		return "", errors.New("Insert: empty data passed")
	}

	dataLen := len(data)
	//fmt.Println(dataLen)
	dataAry := make([]string, dataLen)
	keysAry := make([]string, dataLen)

	schTable, ok := sch.Tables[table]
	if !ok {
		return "", errors.New("Insert: Table map unavailable for table " + table)
	}
	fieldsMap := schTable.Fields
	if fieldsMap == nil {
		return "", errors.New("Insert: Field map unavailable for table " + table)
	}
	i := 0
	for k, v := range data {
		keysAry[i] = k
		//fmt.Println("k=", k, "fieldsMap[k]=", fieldsMap[k], "v=", v)
		r, err := renderInsertValue(fieldsMap[k], v)
		if err != nil {
			return "", err
		}
		dataAry[i] = r
		i++
	}
	fmtStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(keysAry, ","), strings.Join(dataAry, ","))
	return fmtStr, nil
}

func quotedString(value string) string {
	return fmt.Sprintf(`"%s"`, value)
}

func renderBindingInsertValue(f *schema.Field) string {
	return ":" + f.Name
}

func renderInsertValue(f *schema.Field, value interface{}) (string, error) {
	// TODO do we need the schema.Field for more than debugging information?
	switch v := value.(type) {
	case string:
		if v == "" {
			return "", errors.New("renderInsertField: unable to turn the value of " + f.Name + " into string")
		}
		return quotedString(v), nil
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
