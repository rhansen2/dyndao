package sqlitegen

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// TODO: Refactor common code... later. A lot of overall work remains.

// BindingInsert generates the SQL for a given INSERT statement for SQLite with binding parameter values
func (g Generator) BindingInsert(table string, data map[string]interface{}) (string, []interface{}, error) {
	if table == "" {
		return "", nil, errors.New("BindingInsert: Empty table name")
	}
	if data == nil {
		return "", nil, errors.New("BindingInsert: Empty data passed")
	}

	schTable, ok := g.Schema.Tables[table]
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

// BindingUpdate generates the SQL for a given UPDATE statement for SQLite with binding parameter values
func (g Generator) BindingUpdate(sch *schema.Schema, obj *object.Object) (string, []interface{}, []interface{}, error) {
	schTable, ok := g.Schema.Tables[obj.Type]
	if !ok {
		return "", nil, nil, errors.New("BindingUpdate: Table map unavailable for table " + obj.Type)
	}

	fieldsMap := schTable.Fields
	if fieldsMap == nil {
		return "", nil, nil, errors.New("BindingUpdate: Field map unavailable for table " + obj.Type)
	}

	whereClause, bindWhere, err := renderWhereClause(schTable, fieldsMap, obj)
	if err != nil {
		return "", nil, nil, err
	}

	bindArgs := make([]interface{}, len(obj.KV))
	newValuesAry := make([]string, len(obj.KV))
	i := 0
	for k, v := range obj.KV {
		f := fieldsMap[k]
		newValuesAry[i] = fmt.Sprintf("%s = %s%d", f.Name, renderBindingUpdateValue(f), i)
		bindArgs[i] = v
		i++
	}

	// TODO: use schema name from object lookup type, fix in other places...
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s", obj.Type, strings.Join(newValuesAry, ","), whereClause)
	fmt.Println(sqlStr)
	return sqlStr, bindArgs, bindWhere, nil
}

// Insert generates the SQL for a given INSERT statement for SQLite
func (g Generator) Insert(table string, data map[string]interface{}) (string, error) {
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

	schTable, ok := g.Schema.Tables[table]
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
	return fmt.Sprintf(`"%s"`, string(value))
}

func renderBindingInsertValue(f *schema.Field) string {
	return ":" + f.Name
}

func renderBindingUpdateValue(f *schema.Field) string {
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
