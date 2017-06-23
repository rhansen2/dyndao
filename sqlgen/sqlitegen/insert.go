package sqlitegen

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rbastic/dyndao/schema"
)

// Insert generates the SQL for a given INSERT statement for SQLite
func (g Generator) Insert(table string, data map[string]interface{}) (string, error) {
	if table == "" {
		return "", errors.New("Insert: empty table name")
	}
	if data == nil {
		return "", errors.New("Insert: empty data passed")
	}

	dataLen := len(data)
	fmt.Println(dataLen)
	dataAry := make([]string, dataLen)
	keysAry := make([]string, dataLen)
	i := 0
	schTable, ok := g.Schema.Tables[table]
	if !ok {
		return "", errors.New("Table map unavailable for table " + table)
	}
	fieldsMap := schTable.Fields
	if fieldsMap == nil {
		return "", errors.New("Field map unavailable for table " + table)
	}
	for k, v := range data {
		keysAry[i] = k
		fmt.Println("k=", k, "fieldsMap[k]=", fieldsMap[k], "v=", v)
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
