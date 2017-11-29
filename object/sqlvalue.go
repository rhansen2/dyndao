package object

// SQLValue struct is for encapsulating raw SQL Function calls.
// For example, if we want to use SYS_GUID() as a value for an
// INSERT with Oracle, or LAST_INSERT_ID() as a value for an INSERT with MySQL.
// It's meant to be stored in an object's KV, so that it's type
// can be detected and it can be rendered appropriately into a string value.
type SQLValue struct {
	Value string
}

// String method is used to stringify a SQLValue into a raw unquoted
// string.
func (r *SQLValue) String() string {
	return r.Value
}

// NewSQLValue is for creating a new SQLValue. example:
// s := object.NewSQLValue("LAST_INSERT_ID()")
// Passing this as a value for an object using:
// obj := obj.New("SomeDbType")
// obj.Set("primaryKeyColumn", s)
// should result in a query that renders LAST_INSERT_ID() unquoted
// rather than as an unquoted string, so it can be used as a
// function call within SQL.
func NewSQLValue(s string) *SQLValue {
	return &SQLValue{Value: s}
}

// NewNULLValue is syntax sugar for creating a new NULL SQLValue,
// when one wishes to explicitly null something.
func NewNULLValue() *SQLValue {
	return &SQLValue{Value: "NULL"}
}

func (o * Object) ValueIsNULL(v interface{}) bool {
	switch v.(type) {
	case *SQLValue:
		sv := v.(*SQLValue)
		if sv.Value == "NULL" {
			return true
		}
		return false
	default:
		return false
	}
}
