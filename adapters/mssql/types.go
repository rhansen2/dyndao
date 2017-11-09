package mssql

/*
	MSSQL Data types.

	For CLOB support: varchar(max) vs text .. is a bit of a tricky situation. A
string length instead of an integer length (i.e. int(11) vs varchar(max), 11
and max are the 'lengths'.

	NOTE TODO FIXME One thought: perhaps add length as a parameter to the
Is$DATA_TYPE$Type() functions.  Then edge-cases can be dealt with using custom
logic. This may or may not be more suitable than exploring the possibility of
implementing type affinity and data type synonyms.

	The above idea may not be a bad one. Considering Oracle has a NUMBER(i,
j) data type where i and j can determine whether a data type is an integer, a
float, etc., we may need to get even more complex.

	See NOTES file for a link to one of the resources used.
*/

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
	"SMALLINT": true,
	"smallint": true,

	"TINYINT": true,
	"tinyint": true,

	"BIT": true,
	"bit": true,

	// numeric, decimal, money, smallmoney, float, real
}

var floatTypes = map[string]bool{
	"float": true,
	"FLOAT": true,
}

var timestampTypes = map[string]bool{
	// Type affinity?
	"timestamp": true,
	"TIMESTAMP": true,

	// Actual types
	"datetime": true,
	"DATETIME": true,

	// TODO: date, datetime, datetime2, datetimeoffset, smalldatetime, time
}

var lobTypes = map[string]bool{
	"CLOB": true,
	"clob": true,

	"TEXT": true,
	"text": true,

	"NTEXT": true,
	"ntext": true,

	// image is deprecated, varbinary is recommended now
	"image": true,
	"IMAGE": true,
	// TODO: VARCHAR(MAX)
	// TODO: VARBINARY?
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
