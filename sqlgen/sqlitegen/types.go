package sqlitegen

var stringTypes map[string]bool = map[string]bool{
	"TEXT": true,
	"text": true,
}

var numTypes map[string]bool = map[string]bool{
	"INTEGER": true,
	"integer": true,
}

func (g Generator) IsStringType(k string) bool {
	return stringTypes[k]
}

func (g Generator) IsNumberType(k string) bool {
	return numTypes[k]
}
