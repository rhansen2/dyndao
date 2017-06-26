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

	// Note: Name could be 'name', 'email', etc. or 'basic.name', 'basic.email',
	// but it could not be, say, people.basic.name where people is the root namespace
	// and there is also no way to specify the root namespace
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

func fkField(name string) *Field {
	fld := DefaultField()
	fld.Name = name
	fld.DBType = "integer"
	fld.IsIdentity = false
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
	tbl.Fields["PersonID"] = fkField("PersonID")
	tbl.Fields["Zip"] = fieldAddress("Zip")

	tbl.EssentialFields = []string{"AddressID", "PersonID", "Address1", "Address2", "City", "State", "Zip"}

	tbl.ParentTables = []string{"people"}
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
	/*	childTable.ParentTable = "people"
		childTable.LocalField = ""           // TODO: Not needed here?
		childTable.ForeignField = "PersonID" // Field to store our ParentTable record's primary key into
	*/
	personTable.Children["addresses"] = childTable

	addrTable := addressTable()

	sch.Tables["addresses"] = addrTable
	return sch
}
