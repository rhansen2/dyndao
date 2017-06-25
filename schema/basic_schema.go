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

func fieldAddress(n string) *Field {
	fld := DefaultField()
	fld.Source = n
	fld.IsNumber = false
	fld.Name = n
	fld.DBType = "text"
	fld.AllowNull = true // TODO: Should this be used to accept empty keys?
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
	tbl.EssentialFields = []string{"PersonID", "Name"}

	return tbl
}

// Very simple address table initially
func addressTable() *Table {
	tbl := DefaultTable()
	tbl.Primary = "AddressID"
	tbl.MultiKey = true
	tbl.ForeignKeys = []string{"PersonID"}

	tbl.Fields["AddressID"] = primaryField("AddressID")
	tbl.Fields["Address1"] = fieldAddress("Address1")
	tbl.Fields["Address2"] = fieldAddress("Address2")
	tbl.Fields["City"] = fieldAddress("City")
	tbl.Fields["State"] = fieldAddress("State")

	zip := fieldAddress("Zip")
	zip.AllowNull = true // TODO: Is this what we want to use to disable validation?
	tbl.Fields["Zip"] = zip
	return tbl
}

// MockBasicSchema is the basic mock for one table
func MockBasicSchema() *Schema {
	sch := DefaultSchema()
	personTable := peopleTable()

	sch.Tables["people"] = personTable
	return sch
}

// MockNestedSchema is the basic mock for two tables (one table that references a foreign table)
func MockNestedSchema() *Schema {
	sch := MockBasicSchema()
	personTable := sch.Tables["people"]

	childTable := DefaultChildTable()
	childTable.ParentTable = "people"
	childTable.LocalField = ""           // TODO: Not needed here?
	childTable.ForeignField = "PersonID" // Field to store our ParentTable record's primary key into

	personTable.Children["addresses"] = childTable

	sch.Tables["addresses"] = addressTable()
	return sch
}
