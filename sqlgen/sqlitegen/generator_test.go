package sqlitegen

import (
	"fmt"
	"testing"
	"github.com/rbastic/dyndao/schema"

	/*"github.com/rbastic/dyndao/mapper"
	"github.com/rbastic/dyndao/object"*/)

func TestGeneratorBasic(t *testing.T) {
	sch := schema.MockBasicSchema()

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
