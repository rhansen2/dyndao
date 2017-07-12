package oraclegen

var stringTypes = map[string]bool{
	"CLOB":     true,
	"VARCHAR2": true,
	"clob":     true,
	"varchar2": true,
}

var numTypes = map[string]bool{
	"NUMBER": true,
	"number": true,
}

// IsStringType can be used to help determine whether a certain data type is a string type.
// Note that it is case-sensitive.
func (g Generator) IsStringType(k string) bool {
	return stringTypes[k]
}

// IsNumberType can be used to help determine whether a certain data type is a number type.
// Note that it is case-sensitive.
func (g Generator) IsNumberType(k string) bool {
	return numTypes[k]
}

// TODO: strings.ToUpper on key name? just in case?
