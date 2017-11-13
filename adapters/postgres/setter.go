package postgres

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

// TODO: Why was all this necessary for Postgre?

// ripped from https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267
// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 sql.NullInt64

// NullBool is an alias for sql.NullBool data type
type NullBool sql.NullBool

// NullFloat64 is an alias for sql.NullFloat64 data type
type NullFloat64 sql.NullFloat64

// NullTime is an alias for mysql.NullTime data type
//type NullTime mysql.NullTime

// NullString is an alias for sql.NullString data type
type NullString sql.NullString

// Scan implements the Scanner interface for NullString
func (ni *NullString) Scan(value interface{}) error {
	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*ni = NullString{i.String, false}
	} else {
		*ni = NullString{i.String, true}
	}
	return nil
}

// Scan implements the Scanner interface for NullInt64
func (ni *NullInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}
	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*ni = NullInt64{i.Int64, false}
	} else {
		*ni = NullInt64{i.Int64, true}
	}
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
			continue
		} else if s.IsStringType(typeName) {
			val := v.(*NullString)
			// TODO: valid case
			if val.Valid {
				obj.Set(columnNames[i], val.String)
			} else {
				obj.Set(columnNames[i], object.NewNULLValue())
			}
			continue
		} else if s.IsNumberType(typeName) {
			val := v.(*NullInt64)
			if val.Valid {
				obj.Set(columnNames[i], val.Int64)
			} else {
				obj.Set(columnNames[i], object.NewNULLValue())
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
				val := v.(*NullString)
				if val.Valid {
					obj.Set(columnNames[i], val.String)
				}
			} else {
				val := v.(*NullString)
				if val.Valid {
					obj.Set(columnNames[i], val.String)
				}

			}
			continue
		}
		return errors.New("DynamicObjectSetter: Unrecognized type: " + typeName)
	}
	return nil
}

func MakeColumnPointers(s *sg.SQLGenerator, schTable *schema.Table, columnNames []string, columnTypes []*sql.ColumnType) ([]interface{}, error) {
	sliceLen := len(columnNames)
	columnPointers := make([]interface{}, sliceLen)
	for i := 0; i < sliceLen; i++ {
		ct := columnTypes[i]
		typeName := ct.DatabaseTypeName()

		if s.IsNumberType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				var j NullInt64
				columnPointers[i] = &j
			} else {
				var j NullInt64
				columnPointers[i] = &j

			}
		} else if s.IsTimestampType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				// TODO: almost certainly wrong.
				var j time.Time
				columnPointers[i] = &j
			} else {
				var j time.Time
				columnPointers[i] = &j

			}
		} else if s.IsLOBType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				var s NullString
				columnPointers[i] = &s
			} else {
				var s NullString
				columnPointers[i] = &s
			}
		} else if s.IsStringType(typeName) {
			nullable, _ := ct.Nullable()
			if nullable {
				var s NullString
				columnPointers[i] = &s
			} else {
				var s NullString
				columnPointers[i] = &s
			}
		} else {
			return nil, errors.New("MakeColumnPointers: Unrecognized type: " + typeName)
		}
	}
	return columnPointers, nil
}
