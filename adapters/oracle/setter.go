package oracle

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
	"gopkg.in/goracle.v2"
)

// LobDST helps us to implement custom support for goracle.Lob
type LobDST string

// Scan is necessary here to deal with oracle BLOB/CLOB data type.
func (l *LobDST) Scan(src interface{}) error {
	// Scan ignores NULLs, and our DynamicObjectSetter handles them.
	if src == nil {
		return nil
	}

	lob, ok := src.(*goracle.Lob)
	if !ok {
		return fmt.Errorf("LobDST can only be used with goracle.Lib, type was %v", reflect.TypeOf(src))
	}
	res, err := ioutil.ReadAll(lob)
	if err != nil {
		return errors.Wrap(err, "failed to read son")
	}
	*l = LobDST(res)
	return nil
}

// DynamicObjectSetter is used to dynamically set the values of an object by
// checking the necessary types (via sql.ColumnType, and what the driver tells
// us we have for column types)
func DynamicObjectSetter(s *sg.SQLGenerator, schTable *schema.Table, columnNames []string, columnPointers []interface{}, columnTypes []*sql.ColumnType, obj *object.Object) error {
	// NOTE: Read this post for more info on why the code below is written this way:
	// https://stackoverflow.com/questions/23507531/is-golangs-sql-package-incapable-of-ad-hoc-exploratory-queries/23507765#23507765
	for i, v := range columnPointers {
		ct := columnTypes[i]

		typeName := ct.DatabaseTypeName()
		if s.IsTimestampType(typeName) {
			val := v.(*time.Time)
			obj.Set(columnNames[i], *val)
		} else if s.IsStringType(typeName) {
			val := v.(*string)
			obj.Set(columnNames[i], *val)
		} else if s.IsNumberType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				val := v.(*sql.NullInt64)
				if val.Valid {
					obj.Set(columnNames[i], val.Int64)
				}
			} else {
				val := v.(*int64)
				obj.Set(columnNames[i], *val)
			}
		} else if s.IsFloatingType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				val := v.(*sql.NullFloat64)
				if val.Valid {
					obj.Set(columnNames[i], val.Float64)
				}
			} else {
				val := v.(*float64)
				obj.Set(columnNames[i], *val)
			}
		} else if s.IsLOBType(typeName) {
			if v == nil {
				obj.Set(columnNames[i], object.NewNULLValue())
			} else {
				val := v.(*LobDST)
				obj.Set(columnNames[i], string(*val))
			}
		} else {
			return errors.New("dynamicObjectSetter: Unrecognized type: " + typeName)
		}
	}
	return nil
}

func MakeColumnPointers(s *sg.SQLGenerator, schTable *schema.Table, columnNames []string, columnTypes []*sql.ColumnType) ([]interface{}, error) {
	sliceLen := len(columnNames)
	columnPointers := make([]interface{}, sliceLen)
	for i := 0; i < sliceLen; i++ {
		ct := columnTypes[i]
		typeName := ct.DatabaseTypeName()

		if s.IsStringType(typeName) {
			var s string
			columnPointers[i] = &s
		} else if s.IsNumberType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				var j sql.NullInt64
				columnPointers[i] = &j
			} else {
				var j int64
				columnPointers[i] = &j

			}
		} else if s.IsTimestampType(typeName) {
			var j time.Time
			columnPointers[i] = &j
		} else if s.IsLOBType(typeName) {
			s := new(LobDST)
			columnPointers[i] = s
		} else {
			return nil, errors.New("makeColumnPointers: Unrecognized type: " + typeName)
		}
	}
	return columnPointers, nil
}
