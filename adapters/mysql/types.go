package mysql

var stringTypes = map[string]bool{
	"VARSTRING": true,
	"VARCHAR2":  true,
	"varchar2":  true,
	"VARCHAR":   true,
	"varchar":   true,
	"CHAR":      true,
	"char":      true,
}

var numTypes = map[string]bool{
	"NUMBER": true,
	"number": true,
	"INT":    true,
	"int":    true,
}

var floatTypes = map[string]bool{
	"float": true,
	"FLOAT": true,
}

var timestampTypes = map[string]bool{
	"timestamp": true,
	"TIMESTAMP": true,
}

var lobTypes = map[string]bool{
	"TEXT": true,
	"text": true,
	"BLOB": true,
	"blob": true,
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
