package postgres

// TODO: Some of these are unicode types. Do we need to use and support runes instead
// of string here?

var stringTypes = map[string]bool{
	"VARCHAR": true,
	"varchar": true,

	"NVARCHAR": true,
	"nvarchar": true,

	"CHAR": true,
	"char": true,

	"NCHAR": true,
	"nchar": true,
}

var numTypes = map[string]bool{
	"NUMBER": true,
	"number": true,

	"BIGINT": true,
	"bigint": true,

	"INT":      true,
	"int":      true,
	"INT4":     true,
	"int4":     true,
	"SMALLINT": true,
	"smallint": true,

	"TINYINT": true,
	"tinyint": true,

	"BIT": true,
	"bit": true,

	// NOTE: This is for CockroachDB. This was one of the only
	// changes necessary to get it working. Not sure if this
	// originates from an implementation difference or something
	// in the underlying driver.
	"INT8": true,
	"int8": true,

	// numeric, decimal, money, smallmoney, float, real
}

var floatTypes = map[string]bool{
	"float":  true,
	"FLOAT":  true,
	"float8": true,
	"FLOAT8": true,
}

var timestampTypes = map[string]bool{
	"timestamp": true,
	"TIMESTAMP": true,

	// date, datetime, datetime2, datetimeoffset, smalldatetime, time
}

var lobTypes = map[string]bool{
	"TEXT": true,
	"text": true,
}

// IsStringType can be used to help determine whether a certain data type is a string type.
// Note that it is case-sensitive.
func IsStringType(k string) bool {
	return stringTypes[k]
}

// IsNumberType can be used to help determine whether a certain data type is a number type.
// Note that it is case-sensitive.
func IsNumberType(k string) bool {
	return numTypes[k]
}

// IsFloatingType can be used to help determine whether a certain data type is a float type.
// Note that it is case-sensitive.
func IsFloatingType(k string) bool {
	return floatTypes[k]
}

// IsTimestampType can be used to help determine whether a certain data type is a number type.
// Note that it is case-sensitive.
func IsTimestampType(k string) bool {
	return timestampTypes[k]
}

func IsLOBType(k string) bool {
	return lobTypes[k]
}
