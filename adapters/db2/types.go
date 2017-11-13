package db2

var stringTypes = map[string]bool{
	"VARCHAR": true,
	"varchar": true,

	// Deprecated, and don't seem to be working anyway
	// TODO: Flag them somehow to the user?
	// IDEA: Conditional run-time warnings, maybe expand idea of Logger
	// to support them.
	"LONG VARCHAR": true,
	"long varchar": true,

	// Type affinity?
	"TEXT": true,
	"text": true,

	// HACK
	"BLOB": true,
	"blob": true,
}

var numTypes = map[string]bool{
	// TODO: DB2 data types list
	"INTEGER": true,
	"integer": true,

	"NUMBER": true,
	"number": true,
}

var floatTypes = map[string]bool{
	// TODO: DB2 data types list
	"float": true,
	"FLOAT": true,
}

var timestampTypes = map[string]bool{
	// TODO: DB2 data types list
	"timestamp": true,
	"TIMESTAMP": true,
}

var lobTypes = map[string]bool{
	// FIXME: Mixing of CLOB and (VAR)BINARY types, ergh
	"CLOB":   true,
	"clob":   true,
	"BINARY": true,
	"binary": true,

	"VARBINARY": true,
	"varbinary": true,
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
