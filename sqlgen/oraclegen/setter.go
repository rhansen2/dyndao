package oraclegen

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/object"
	"gopkg.in/goracle.v2"
	"io/ioutil"
	"time"
)

// LobDST helps us to implement custom support for goracle.Lob
type LobDST string

// Scan is necessary here to deal with oracle BLOB/CLOB data type.
func (l *LobDST) Scan(src interface{}) error {
	lob, ok := src.(*goracle.Lob)
	if !ok {
		return errors.New("LobDST can only be used with goracle.Lib")
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
func (s Generator) DynamicObjectSetter(columnNames []string, columnPointers []interface{}, columnTypes []*sql.ColumnType, obj *object.Object) error {
	// NOTE: Read this post for more info on why the code below is written this way:
	// https://stackoverflow.com/questions/23507531/is-golangs-sql-package-incapable-of-ad-hoc-exploratory-queries/23507765#23507765
	for i, v := range columnPointers {
		ct := columnTypes[i]

		typeName := ct.DatabaseTypeName()
		// TODO: Not sure this is actually correct?
		if s.IsTimestampType(typeName) {
			val := v.(*time.Time)
			obj.Set(columnNames[i], *val)
		} else if s.IsStringType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				// TODO: Does this work properly across databases?
				//val := v.(*sql.NullString)
				val := v.(*string)
				obj.Set(columnNames[i], *val)
				/*
					if val.Valid {
					}
				*/
				// TODO: We don't set keys for null values. How else can we support this?
			} else {
				val := v.(*string)
				obj.Set(columnNames[i], *val)

			}
		} else if s.IsNumberType(typeName) {
			// TODO: support more than 'int64' for integer...?
			nullable, _ := ct.Nullable()
			if nullable {
				val := v.(*sql.NullInt64)
				if val.Valid {
					obj.Set(columnNames[i], val.Int64)
				}
				// TODO: We don't set keys for null values. How else can we support this?
			} else {
				val := v.(*int64)
				obj.Set(columnNames[i], *val)
			}
		} else if s.IsFloatingType(typeName) {
			// TODO: support more than 'int64' for integer...?
			nullable, _ := ct.Nullable()
			if nullable {
				val := v.(*sql.NullFloat64)
				if val.Valid {
					obj.Set(columnNames[i], val.Float64)
				}
				// TODO: We don't set keys for null values. How else can we support this?
			} else {
				val := v.(*float64)
				obj.Set(columnNames[i], *val)
			}
		} else if s.IsLOBType(typeName) {
			val := v.(*LobDST)
			obj.Set(columnNames[i], string(*val))
		} else {
			return errors.New("dynamicObjectSetter: Unrecognized type: " + typeName)
		}
		// TODO: add timestamp support.?
	}
	return nil
}

func (s Generator) MakeColumnPointers(sliceLen int, columnTypes []*sql.ColumnType) ([]interface{}, error) {
	columnPointers := make([]interface{}, sliceLen)
	for i := 0; i < sliceLen; i++ {
		ct := columnTypes[i]
		typeName := ct.DatabaseTypeName()
		if s.IsStringType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				var s string
				columnPointers[i] = &s
			} else {
				var s string
				columnPointers[i] = &s
			}
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
				s := new(LobDST)
				columnPointers[i] = s
			} else {
				s := new(LobDST)
				columnPointers[i] = s
			}
		} else {
			return nil, errors.New("makeColumnPointers: Unrecognized type: " + typeName)
		}
	}
	return columnPointers, nil
}
