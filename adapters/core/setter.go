package core

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/object"
	sg "github.com/rbastic/dyndao/sqlgen"
	"time"
)

// DynamicObjectSetter is used to dynamically set the values of an object by
// checking the necessary types (via sql.ColumnType, and what the driver tells
// us we have for column types)
func DynamicObjectSetter(s *sg.SQLGenerator, columnNames []string, columnPointers []interface{}, columnTypes []*sql.ColumnType, obj *object.Object) error {
	// NOTE: Read this post for more info on why the code below is written this way:
	// https://stackoverflow.com/questions/23507531/is-golangs-sql-package-incapable-of-ad-hoc-exploratory-queries/23507765#23507765
	for i, v := range columnPointers {
		ct := columnTypes[i]

		typeName := ct.DatabaseTypeName()

		if s.IsTimestampType(typeName) {
			val := v.(*time.Time)
			obj.Set(columnNames[i], *val)
			continue
		} else if s.IsStringType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				val := v.(*sql.NullString)
				obj.Set(columnNames[i], *val)
			} else {
				val := v.(*string)
				obj.Set(columnNames[i], *val)

			}
			continue
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
			continue
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
			continue
		} else if s.IsLOBType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				val := v.(*sql.NullString)
				obj.Set(columnNames[i], *val)
			} else {
				val := v.(*string)
				obj.Set(columnNames[i], *val)

			}
			continue
		}
		return errors.New("DynamicObjectSetter: Unrecognized type: " + typeName)
	}
	return nil
}

func MakeColumnPointers(s *sg.SQLGenerator, sliceLen int, columnTypes []*sql.ColumnType) ([]interface{}, error) {
	columnPointers := make([]interface{}, sliceLen)
	for i := 0; i < sliceLen; i++ {
		ct := columnTypes[i]
		typeName := ct.DatabaseTypeName()

		if typeName == "" {
			panic("dyndao MakeColumnPointers: ct.DatabaseTypeName() does not appear to be implemented - typeName was an empty string")
		}

		if s.IsNumberType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				var j sql.NullInt64
				columnPointers[i] = &j
			} else {
				var j int64
				columnPointers[i] = &j

			}
		} else if s.IsTimestampType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				var j time.Time
				columnPointers[i] = &j
			} else {
				var j time.Time
				columnPointers[i] = &j

			}
		} else if s.IsLOBType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				var s sql.NullString
				columnPointers[i] = &s
			} else {
				var s string
				columnPointers[i] = &s
			}
		} else if s.IsStringType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				var s sql.NullString
				columnPointers[i] = &s
			} else {
				var s string
				columnPointers[i] = &s
			}
		} else {
			return nil, errors.New("MakeColumnPointers: Unrecognized type: " + typeName)
		}
	}
	return columnPointers, nil
}
