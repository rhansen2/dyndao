package test

import (
	"github.com/rbastic/dyndao/schema"
)

// Basic test mock
func fieldName() *schema.Column {
	fld := schema.DefaultColumn()
	fld.IsNumber = false
	fld.Name = "Name"
	fld.DBType = "text"
	fld.Source = "Name"
	//	fld.Length = 36
	return fld
}

func fieldMetadata() *schema.Column {
	fld := schema.DefaultColumn()
	fld.IsNumber = false
	fld.Name = "NullableMetadata"
	fld.DBType = "text"
	fld.Source = "NullableMetadata"
	fld.AllowNull = true
	return fld
}

func fieldAddress(n string) *schema.Column {
	fld := schema.DefaultColumn()
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

func primaryColumn(name string) *schema.Column {
	fld := schema.DefaultColumn()
	fld.Name = name
	fld.DBType = "integer"
	fld.IsIdentity = true
	fld.IsNumber = true
	fld.Source = name
	return fld
}

// Very simple person table initially
func peopleTable() *schema.Table {
	tbl := schema.DefaultTable()
	tbl.Name = "people"
	tbl.MultiKey = false
	tbl.Primary = "PersonID"
	tbl.Columns["PersonID"] = primaryColumn("PersonID")
	tbl.Columns["Name"] = fieldName()
	tbl.Columns["NullableMetadata"] = fieldMetadata()
	tbl.EssentialColumns = []string{"PersonID", "Name"}

	return tbl
}

func fkColumn(name string) *schema.Column {
	fld := schema.DefaultColumn()
	fld.Name = name
	fld.DBType = "integer"
	fld.IsIdentity = false
	fld.IsNumber = true
	fld.Source = name
	return fld
}

// Very simple address table initially
func addressTable() *schema.Table {
	tbl := schema.DefaultTable()
	tbl.Name = "addresses"
	tbl.Primary = "AddressID"
	tbl.MultiKey = true
	tbl.ForeignKeys = []string{"PersonID"}

	tbl.Columns["AddressID"] = primaryColumn("AddressID")
	tbl.Columns["Address1"] = fieldAddress("Address1")
	tbl.Columns["Address2"] = fieldAddress("Address2")
	tbl.Columns["City"] = fieldAddress("City")
	tbl.Columns["State"] = fieldAddress("State")
	tbl.Columns["PersonID"] = fkColumn("PersonID")
	tbl.Columns["Zip"] = fieldAddress("Zip")

	tbl.EssentialColumns = []string{"AddressID", "PersonID", "Address1", "Address2", "City", "State", "Zip"}

	tbl.ParentTables = []string{"people"}
	return tbl
}

// MockBasicSchema is the basic mock for one table
func MockBasicSchema() *schema.Schema {
	sch := schema.DefaultSchema()
	personTable := peopleTable()

	sch.Tables["people"] = personTable
	return sch
}

// MockNestedSchema is the basic mock for two tables (one table that references a foreign table)
func MockNestedSchema() *schema.Schema {
	sch := MockBasicSchema()
	personTable := sch.Tables["people"]

	childTable := schema.DefaultChildTable()
	personTable.Children["addresses"] = childTable

	addrTable := addressTable()
	sch.Tables["addresses"] = addrTable
	return sch
}
