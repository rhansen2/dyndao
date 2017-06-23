package schema

// Basic test mock
func fieldName() *Field {
	fld := DefaultField()
	fld.IsNumber = false
	fld.Name = "Name"
	fld.DBType = "text"
	fld.Source = "Name"
	//	fld.Length = 36
	return fld
}

func primaryField(name string) *Field {
	fld := DefaultField()
	fld.Name = name
	fld.DBType = "integer"
	fld.IsIdentity = true
	fld.IsNumber = true
	fld.Source = name
	return fld
}

// Very simple person table initially
func peopleTable() *Table {
	tbl := DefaultTable()
	tbl.MultiKey = false
	tbl.Primary = "PersonID"
	tbl.Fields["PersonID"] = primaryField("PersonID") //fieldID("PersonID")
	tbl.Fields["Name"] = fieldName()
	return tbl
}

// Basic mock  one table
func MockBasicSchema() *Schema {
	sch := DefaultSchema()
	sch.Tables["people"] = peopleTable()
	return sch
}
