package sqlitegen

import (
	"fmt"
	"testing"
	"github.com/rbastic/dyndao/schema"

	/*"github.com/rbastic/dyndao/mapper"
	"github.com/rbastic/dyndao/object"*/
)

func TestGeneratorBasic(t *testing.T) {
	sch := schema.MockBasicSchema()

	// Basic generator initialization
	basic := New("sqlite", "testDB", sch)

	sql, err := basic.CreateTable("people")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql + ";")
	// TODO: test output or execution
	// TODO: Run schema against sqlite3

	fmt.Println(basic.Begin("") + ";")
	// TODO: test output or execution

	sql, err = basic.Insert(
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
	basic := New("sqlite", "testDB", sch)

	sql, err := basic.CreateTable("people")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql + ";")

	sql, err = basic.CreateTable("addresses")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql + ";")
	// TODO: test output or execution
	// TODO: Run schema against sqlite3

	fmt.Println(basic.Begin("") + ";")
	// TODO: test output or execution

	sql, err = basic.Insert(
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
