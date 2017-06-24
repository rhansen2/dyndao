package schema

// DefaultSchema returns an empty schema ready to be populated
func DefaultSchema() *Schema {
	tables := make(map[string]*Table)
	sch := &Schema{Tables: tables}

	return sch
}

// DefaultTable returns an empty table ready to be populated
func DefaultTable() *Table {
	fieldsMap := make(map[string]*Field)
	childrenMap := make(map[string]*ChildTable)

	tbl := &Table{
		MultiKey:        false,
		Primary:         "",
		Primaries:       nil,
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
		Source: "",
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
