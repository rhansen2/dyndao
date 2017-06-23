package sqlitegen

import (
	"fmt"
	"testing"

	"github.com/rbastic/dyndao/schema"
	/*"github.com/rbastic/dyndao/mapper"
	"github.com/rbastic/dyndao/object"*/)

// Basic schema test mock
func fieldName() *schema.Field {
	fld := schema.DefaultField()
	fld.IsNumber = false
	fld.Name = "Name"
	fld.DBType = "text"
	//	fld.Length = 36
	return fld
}

func fieldID(name string) *schema.Field {
	fld := schema.DefaultField()
	fld.Name = name
	fld.IsNumber = true
	return fld
}

func primaryField(name string) *schema.Field {
	fld := schema.DefaultField()
	fld.Name = name
	fld.DBType = "integer"
	fld.IsIdentity = true
	fld.IsNumber = true
	return fld
}

// Very simple person table initially
func peopleTable() *schema.Table {
	tbl := schema.DefaultTable()
	tbl.MultiKey = false
	tbl.Primary = "PersonID"
	tbl.Fields["PersonID"] = primaryField("PersonID") //fieldID("PersonID")
	tbl.Fields["Name"] = fieldName()
	return tbl
}

// Basic mock schema, one table
func BasicSchema() *schema.Schema {
	sch := schema.DefaultSchema()
	sch.Tables["people"] = peopleTable()
	return sch
}

func TestGeneratorBasic(t *testing.T) {
	sch := BasicSchema()

	// Basic generator initialization
	basic := New("sqlite", "testDB", sch)

	fmt.Println(basic.CreateTable("people"))
	// TODO: Run schema against sqlite3

	fmt.Println(basic.Begin(""))

	fmt.Println(basic.Insert(
		"people",
		map[string]interface{}{
			"PersonID": 1,
			"Name":     "Sam",
		},
	))

	fmt.Println(basic.Commit())

	fmt.Println(basic.DropTable("people"))
}
