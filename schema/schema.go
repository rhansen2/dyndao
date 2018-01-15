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

// GetColumnName returns the correct field name in a potentially aliased environment.
// This is useful in situations where you aren't sure what the 'real' key name
// may potentially be.
func (t *Table) GetColumnName(n string) string {
	if t.ColumnAliases != nil {
		realName, ok := t.ColumnAliases[n]
		if !ok {
			// Perhaps it is not an alias
			return n
		}
		return realName
	}
	return n
}

// GetColumn returns the correct field in a potentially aliased environment.
func (t *Table) GetColumn(n string) *Column {
	if t.ColumnAliases != nil {
		realName, ok := t.ColumnAliases[n]
		if !ok {
			// Perhaps it is not an alias
			return t.Columns[n]
		}
		return t.Columns[realName]
	}
	return t.Columns[n]
}

// DefaultTable returns an empty table ready to be populated
func DefaultTable() *Table {
	fieldsMap := make(map[string]*Column)
	childrenMap := make(map[string]*ChildTable)
	emptyAliasesMap := make(map[string]string)

	tbl := &Table{
		MultiKey:         false,
		Primary:          "",
		Columns:          fieldsMap,
		EssentialColumns: nil,
		Children:         childrenMap,
		ColumnAliases:    emptyAliasesMap,
	}
	return tbl
}

// DefaultColumn returns an empty field struct ready to be populated
func DefaultColumn() *Column {
	fld := &Column{
		Name:         "",
		AllowNull:    false,
		DefaultValue: "",
		IsNumber:     false,
		DBType:       "",
		IsIdentity:   false,
	}
	return fld
}

// DefaultChildTable returns an empty child table ready to be populated
func DefaultChildTable() *ChildTable {
	chld := &ChildTable{
		ParentTable:   "",
		MultiKey:      false,
		LocalColumn:   "",
		ForeignColumn: "",

		LocalColumns:   nil,
		ForeignColumns: nil,
	}
	return chld
}
