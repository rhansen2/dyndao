package sqlitegen

// Derived from: www.sqlite.org/datatype3.html
// TODO: That table only shows a small subset of the datatypes SQLite will accept.
var stringTypes = map[string]bool{
	"TEXT":              true,
	"text":              true,
	"CHARACTER":         true,
	"character":         true,
	"VARCHAR":           true,
	"varchar":           true,
	"VARYING CHARACTER": true,
	"varying character": true,
	"NCHAR":             true,
	"nchar":             true,
	"NATIVE CHARACTER":  true,
	"native character":  true,
	"NVARCHAR":          true,
	"nvarchar":          true,
	"CLOB":              true,
	"clob":              true,
}

var numTypes = map[string]bool{
	"INTEGER":          true,
	"integer":          true,
	"INT":              true,
	"int":              true,
	"TINYINT":          true,
	"tinyint":          true,
	"SMALLINT":         true,
	"smallint":         true,
	"MEDIUMINT":        true,
	"mediumint":        true,
	"BIGINT":           true,
	"bigint":           true,
	"UNSIGNED BIG INT": true,
	"unsigned big int": true,
	"INT2":             true,
	"int2":             true,
	"INT8":             true,
	"int8":             true,
}

var floatTypes = map[string]bool{
	"REAL":             true,
	"real":             true,
	"DOUBLE":           true,
	"double":           true,
	"DOUBLE PRECISION": true,
	"double precision": true,
	"FLOAT":            true,
	"float":            true,
}

var timestampTypes = map[string]bool{
	"datetime": true,
	"DATETIME": true,
}

// TODO: blob types

// IsStringType can be used to help determine whether a certain data type is a string type.
func (g Generator) IsStringType(k string) bool {
	return stringTypes[k]
}

// IsFloatingType can be used to help determine whether a certain data type is a floating point type.
func (g Generator) IsFloatingType(k string) bool {
	return floatTypes[k]
}

// IsNumberType can be used to help determine whether a certain data type is a number type.
func (g Generator) IsNumberType(k string) bool {
	return numTypes[k]
}

// TODO: strings.ToUpper on key name? just in case?

// IsTimestampType can be used to help determine whether a certain data type is a number type.
// Note that it is case-sensitive.
func (g Generator) IsTimestampType(k string) bool {
	return timestampTypes[k]
}

// IsLOBType remains unimplemented for SQLite.
func (g Generator) IsLOBType(k string) bool {
	return false
}
