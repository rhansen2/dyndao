package schema

import (
	"fmt"
	"testing"
)

func TestDefaultSchema(t *testing.T) {
	sch := DefaultSchema()
	fmt.Println(sch)
}

func TestDefaultTable(t *testing.T) {
	tbl := DefaultTable()
	fmt.Println(tbl)
}

func TestDefaultField(t *testing.T) {
	fld := DefaultField()
	fmt.Println(fld)
}

func TestDefaultChildTable(t *testing.T) {
	chld := DefaultChildTable()
	fmt.Println(chld)
}

func fieldName() *Field {
	fld := DefaultField()
	fld.IsNumber = false
	fld.Name = "Name"
	fld.DBType = "varchar"
	fld.Length = 36
	return fld
}

func fieldID(name string) *Field {
	fld := DefaultField()
	fld.Name = name
	fld.IsNumber = true
	return fld
}

func peopleTable() *Table {
	tbl := DefaultTable()
	tbl.MultiKey = false
	tbl.Primary = "PersonID"
	tbl.Fields["PersonID"] = fieldID("PersonID")
	tbl.Fields["Name"] = fieldName()
	return tbl
}

func jobTable() *Table {
	tbl := DefaultTable()
	tbl.MultiKey = false
	tbl.Primary = "JobID"
	tbl.Fields["JobID"] = fieldID("JobID")
	tbl.Fields["Name"] = fieldName()
	return tbl
}
func permissionsTable() *Table {
	tbl := DefaultTable()
	tbl.MultiKey = false
	//	tbl.Primary = "PermissionsID"
	tbl.Fields["PermissionsID"] = fieldID("PermissionsID")
	tbl.Fields["Name"] = fieldName()
	return tbl
}

func TestSchemaBasic(t *testing.T) {
	_ = basicSchema()
}

func basicSchema() *Schema {
	sch := DefaultSchema()

	sch.Tables["people"] = peopleTable()
	sch.Tables["job"] = jobTable()
	sch.Tables["permissions"] = permissionsTable()
	return sch
}

/*
func usersTable() *Schema {
	tbl := DefaultTable()
	tbl.MultiKey = false
	tbl.Primary = "id"
	tbl.Fields[""]
}

func addressesTable() *Schema {

}

func addressBookSchema() *Schema {

	sch := DefaultSchema()
	sch.Tables["users"] = usersTable()
	sch.Tables["addresses"] = addressesTable()
}

*/
