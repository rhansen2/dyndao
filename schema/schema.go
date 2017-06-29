package schema

import (
	"encoding/json"
)

// DefaultSchema returns an empty schema ready to be populated
func DefaultSchema() *Schema {
	tables := make(map[string]*Table)
	sch := &Schema{Tables: tables}

	return sch
}

func (s *Schema) ToJSON() (string, error) {
	buf, err := s.ToJSONBytes()
	if err != nil {
		return "", err
	}
	return string(buf), err
}

func (s *Schema) ToJSONBytes() ([]byte, error) {
	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func FromJSON(jsonStr string) (*Schema, error) {
	sch := DefaultSchema()
	err := json.Unmarshal([]byte(jsonStr), &sch)
	if err != nil {
		return nil, err
	}
	return sch, nil
}

func FromJSONBytes(jsonBytes []byte) (*Schema, error) {
	sch := DefaultSchema()
	err := json.Unmarshal([]byte(jsonBytes), &sch)
	if err != nil {
		return nil, err
	}
	return sch, nil
}

// DefaultTable returns an empty table ready to be populated
func DefaultTable() *Table {
	fieldsMap := make(map[string]*Field)
	childrenMap := make(map[string]*ChildTable)

	tbl := &Table{
		MultiKey:        false,
		Primary:         "",
		Fields:          fieldsMap,
		EssentialFields: nil,
		Children:        childrenMap,
	}
	return tbl
}

// DefaultField returns an empty field struct ready to be populated
func DefaultField() *Field {
	fld := &Field{
		Name:         "",
		AllowNull:    false,
		DefaultValue: "",
		IsNumber:     false,
		DBType:       "",
		IsIdentity:   false,
		Source:       "",
	}
	return fld
}

// DefaultChildTable returns an empty child table ready to be populated
func DefaultChildTable() *ChildTable {
	chld := &ChildTable{
		ParentTable:  "",
		MultiKey:     false,
		LocalField:   "",
		ForeignField: "",

		LocalFields:   nil,
		ForeignFields: nil,
	}
	return chld
}
