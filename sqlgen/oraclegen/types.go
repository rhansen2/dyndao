package oraclegen

var stringTypes map[string]bool = map[string]bool{
	"CLOB":     true,
	"VARCHAR2": true,
	"clob":     true,
	"varchar2": true,
}

var numTypes map[string]bool = map[string]bool{
	"NUMBER": true,
	"number": true,
}

func (g Generator) IsStringType(k string) bool {
	return stringTypes[k]
}

func (g Generator) IsNumberType(k string) bool {
	return numTypes[k]
}
