package schema

import (
	"encoding/json"
	//"fmt"
)

// DefaultSchema returns an empty schema ready to be populated
func DefaultSchema() *Schema {
	tables := make(map[string]*Table)
	sch := &Schema{Tables: tables, TableAliases: nil}

	return sch
}

// ToJSON converts a schema into a JSON string.
func (s *Schema) ToJSON() (string, error) {
	buf, err := s.ToJSONBytes()
	if err != nil {
		return "", err
	}
	return string(buf), err
}

// ToJSONBytes converts a schema into a JSON byte array.
func (s *Schema) ToJSONBytes() ([]byte, error) {
	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// FromJSON unmarshals a JSON string into a Schema object.
func FromJSON(jsonStr string) (*Schema, error) {
	sch := DefaultSchema()
	err := json.Unmarshal([]byte(jsonStr), &sch)
	if err != nil {
		return nil, err
	}
	return sch, nil
}

// FromJSONBytes unmarshals a JSON byte array into a Schema object.
func FromJSONBytes(jsonBytes []byte) (*Schema, error) {
	sch := DefaultSchema()
	err := json.Unmarshal(jsonBytes, &sch)
	if err != nil {
		return nil, err
	}
	return sch, nil
}

// GetTableName returns the correct Table name in a potentially aliased environment.
func (s *Schema) GetTableName(n string) string {
	if s.TableAliases != nil {
		realName, ok := s.TableAliases[n]
		if !ok {
			// Perhaps it is not an alias
			return n
		}
		return realName
	}
	return n
}

// GetTable returns the correct Table type in a potentially aliased environment.
func (s *Schema) GetTable(n string) *Table {
	if s.TableAliases != nil {
		realName, ok := s.TableAliases[n]
		if !ok {
			// Perhaps it is not an alias
			return s.Tables[n]
		}
		return s.Tables[realName]
	}
	return s.Tables[n]
}

// GetFieldName returns the correct field name in a potentially aliased environment.
// This is useful in situations where you aren't sure what the 'real' key name
// may potentially be.
func (t *Table) GetFieldName(n string) string {
	if t.FieldAliases != nil {
		realName, ok := t.FieldAliases[n]
		if !ok {
			// Perhaps it is not an alias
			return n
		}
		return realName
	}
	return n
}

// GetField returns the correct field in a potentially aliased environment.
func (t *Table) GetField(n string) *Field {
	if t.FieldAliases != nil {
		realName, ok := t.FieldAliases[n]
		if !ok {
			// Perhaps it is not an alias
			return t.Fields[n]
		}
		return t.Fields[realName]
	}
	return t.Fields[n]
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
