package mysqlgen

import (
	"fmt"
	"testing"

	"github.com/rbastic/dyndao/schema"
)

func TestGeneratorBasic(t *testing.T) {
	sch := schema.MockBasicSchema()
	basic := New("testDB", sch) // Basic generator initialization

	sql, err := basic.CreateTable(sch, "people")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql + ";")
	fmt.Println(basic.Begin("") + ";")

	sqlStr, bindArgs, err := basic.BindingInsert(
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
	fmt.Println(sqlStr + ";")
	fmt.Println("bindArgs:", bindArgs)
	fmt.Println(basic.Commit() + ";")

	fmt.Println(basic.DropTable("people") + ";")
}

func TestGeneratorNested(t *testing.T) {
	sch := schema.MockNestedSchema()

	// Basic generator initialization
	basic := New("testDB", sch)

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

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql + ";")
	fmt.Println(basic.Commit() + ";")
	fmt.Println(basic.DropTable("people") + ";")
}
