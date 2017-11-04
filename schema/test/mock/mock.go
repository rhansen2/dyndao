package mock

import (
	"github.com/rbastic/dyndao/schema"
)

// Basic test mock
func fieldName() *schema.Column {
	fld := schema.DefaultColumn()
	fld.IsNumber = false
	fld.Name = "Name"
	fld.DBType = "text"
	return fld
}

func fieldNullText() *schema.Column {
	fld := schema.DefaultColumn()
	fld.IsNumber = false
	fld.Name = "NullText"
	fld.DBType = "text"
	fld.AllowNull = true
	return fld
}

func fieldNullInt() *schema.Column {
	fld := schema.DefaultColumn()
	fld.IsNumber = true
	fld.Name = "NullInt"
	fld.DBType = "integer"
	fld.AllowNull = true
	return fld
}

// TODO: Int field

func fieldNullVarchar() *schema.Column {
	fld := schema.DefaultColumn()
	fld.IsNumber = true
	fld.Name = "NullVarchar"
	fld.DBType = "varchar"
	fld.AllowNull = true
	fld.Length = 30
	return fld
}

func fieldNullBlob() *schema.Column {
	fld := schema.DefaultColumn()
	fld.IsNumber = true
	fld.Name = "NullBlob"
	fld.DBType = "blob"
	fld.AllowNull = true
	return fld
}

func fieldAddress(n string) *schema.Column {
	fld := schema.DefaultColumn()
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
	tbl.Columns["NullText"] = fieldNullText()
	tbl.Columns["NullInt"] = fieldNullInt()
	tbl.Columns["NullVarchar"] = fieldNullVarchar()
	tbl.Columns["NullBlob"] = fieldNullBlob()

	// TODO: Why was NullText getting retrieved as NULL when we didn't
	// have it in the EssentialColumns list?
	tbl.EssentialColumns = []string{"PersonID", "Name", "NullText", "NullInt", "NullVarchar", "NullBlob"}

	return tbl
}

func fkColumn(name string) *schema.Column {
	fld := schema.DefaultColumn()
	fld.Name = name
	fld.DBType = "integer"
	fld.IsIdentity = false
	fld.IsNumber = true
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

// BasicSchema is the basic mock for one table
func BasicSchema() *schema.Schema {
	sch := schema.DefaultSchema()
	personTable := peopleTable()

	sch.Tables["people"] = personTable
	return sch
}

// NestedSchema is the basic mock for two tables (one table that references a foreign table)
func NestedSchema() *schema.Schema {
	sch := BasicSchema()
	personTable := sch.Tables["people"]

	childTable := schema.DefaultChildTable()
	personTable.Children["addresses"] = childTable

	addrTable := addressTable()
	sch.Tables["addresses"] = addrTable
	return sch
}
