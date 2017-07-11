package sqlitegen

import (
	"fmt"
	"testing"

	"github.com/rbastic/dyndao/schema"
)

func TestGeneratorBasic(t *testing.T) {
	sch := schema.MockBasicSchema()
	basic := New("testDB", sch, false) // Basic generator initialization

	sql, err := basic.CreateTable(sch, "people")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql + ";")
	// TODO: test output or execution
	// TODO: Run schema against sqlite3

	fmt.Println(basic.Begin("") + ";")
	// TODO: test output or execution

	sql, err = basic.Insert(
		sch,
		"people",
		map[string]interface{}{
			"PersonID": 1,
			"Name":     "Sam",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql + ";")

	// TODO: test output or execution
	fmt.Println(basic.Commit() + ";")

	fmt.Println(basic.DropTable("people") + ";")
	// TODO: test output or execution
}

func TestGeneratorNested(t *testing.T) {
	sch := schema.MockNestedSchema()

	// Basic generator initialization
	basic := New("testDB", sch, false)

	sql, err := basic.CreateTable(sch, "people")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql + ";")

	sql, err = basic.CreateTable(sch, "addresses")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql + ";")

	fmt.Println(basic.Begin("") + ";")

	sql, err = basic.Insert(
		sch,
		"people",
		map[string]interface{}{
			"PersonID": 1,
			"Name":     "Sam",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql + ";")

	fmt.Println(basic.Commit() + ";")

	fmt.Println(basic.DropTable("people") + ";")
}
