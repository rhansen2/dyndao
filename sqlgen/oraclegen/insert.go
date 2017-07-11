package oraclegen

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// TODO: Refactor common code... later. A lot of overall work remains.

func (g Generator) setPKifNeeded(data map[string]interface{}, identityCol string) {
	if g.CallerSuppliesPK {
		// NOTE: Currently we support this as a SYS_GUID implementation,
		// and technically, /we/ aren't the caller ... maybe I will rename
		// this as ORMSuppliesPK ... but I don't have another reference
		// point to compare against yet. YAGNI applies.
		_, ok := data[identityCol]
		if !ok {
			data[identityCol] = object.NewSqlValue("SYS_GUID()")
		}
	}
}

func (g Generator) coreBindingInsert(data map[string]interface{}, identityCol string, fieldsMap map[string]*schema.Field) ([]string, []string, []interface{}) {
	dataLen := len(data)
	bindNames := make([]string, dataLen)
	colNames := make([]string, dataLen)
	bindArgs := make([]interface{}, dataLen)
	i := 0
	for k, v := range data {
		colNames[i] = fmt.Sprintf(`%s`, k)
		var r string

		if g.CallerSuppliesPK && k == identityCol {
			switch typ := v.(type) {
			case int64:
				r = string(v.(int64))
			case string:
				r = v.(string)
			case *object.SqlValue:
				thing := v.(*object.SqlValue)
				r = thing.Value
				v = nil
			default:
				panic(fmt.Sprintf("Unknown type [%v] in switch", typ))
			}
		}
		if r == "" {
			r = renderBindingInsertValue(fieldsMap[k])
		}
		bindNames[i] = fmt.Sprintf(`%s`, r)
		if v == nil {
			bindArgs[i] = v
		} else {
			barg, err := renderInsertValue(fieldsMap[k], v)
			if err != nil {
				panic(err.Error())
			}
			bindArgs[i] = barg

		}
		fmt.Printf("i=%v bindNames[i]=%v bindArgs[i]=%v\n", i, bindNames[i], bindArgs[i])
		i++
	}
	return bindNames, colNames, bindArgs
}

// TODO: Push these 3 routines onto github.
func countNils(maybeNils []interface{}) int {
	var count int
	for _, v := range maybeNils {
		if v == nil {
			count++
		}
	}
	return count
}

func removeNils(someNils []interface{}, count int) []interface{} {
	noNils := make([]interface{}, len(someNils)-count)
	j := 0
	for _, v := range someNils {
		if v != nil {
			noNils[j] = v
			j++
		}

	}
	return noNils
}

func removeNilsIfNeeded(maybeNils []interface{}) []interface{} {
	numNils := countNils(maybeNils)
	if numNils > 0 {
		return removeNils(maybeNils, numNils)
	}
	return maybeNils
}

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

	tableName := schema.GetTableName(schTable.Name, table)

	fieldsMap := schTable.Fields
	if fieldsMap == nil {
		return "", nil, errors.New("BindingInsert: Field map unavailable for table " + table)
	}

	identityCol := schTable.Primary
	g.setPKifNeeded(data, identityCol)

	bindNames, colNames, bindArgs := g.coreBindingInsert(data, identityCol, fieldsMap)
	bindArgs = removeNilsIfNeeded(bindArgs)
	fmt.Println("bindArgs->", bindArgs)
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s /*LASTINSERTID*/ INTO :%s",
		tableName,
		strings.Join(colNames, ","),
		strings.Join(bindNames, ","),
		identityCol,
		identityCol)

	fmt.Println("DEBUG: INSERT sqlStr->", sqlStr, "bindArgs->", bindArgs)
	return sqlStr, bindArgs, nil
}

func quotedString(value string) string {
	// TODO: Quote the data according to semantics of local database. This should be provided by dyndao's sql generator?
	// ( TODO: Research what something like xorm might provide already in this area... )
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
	switch typ := value.(type) {
	case string:
		str := value.(string)
		if str == "" {
			return "", errors.New("renderInsertField: unable to turn the value of " + f.Name + " into string")
		}
		return str, nil
	case int32:
		num := value.(int32)
		return string(num), nil
	case int:
		num := value.(int)
		return strconv.Itoa(num), nil
	case int64:
		num := value.(int64)
		return fmt.Sprintf("%d", num), nil
	case uint64:
		num := value.(uint64)
		return fmt.Sprintf("%d", num), nil
	case float64:
		num := value.(float64)
		if f.IsNumber {
			return fmt.Sprintf("%d", int64(num)), nil
		}
		// TODO: when we support more than regular integers, we'll need to care about this more
		return fmt.Sprintf("%f", num), nil
	case *object.SqlValue:
		val := value.(*object.SqlValue)
		return val.String(), nil
	case object.SqlValue:
		val := value.(object.SqlValue)
		return val.String(), nil
	default:
		return "", fmt.Errorf("renderInsertField: unknown type %v for the value of %s", typ, f.Name)
	}
}

// TODO: InsertBinding
