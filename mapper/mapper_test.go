package mapper

import (
	"fmt"
	"testing"
	"github.com/rbastic/dyndao/schema"
)

func TestBasicMapper(t *testing.T) {
	sch := schema.MockBasicSchema()

	obj, err := toObjectFromJSON( sch, "people", getJsonData() )

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(obj)
}

func getJsonData() string {
	return `
{
	"PersonID": 1,
	"Name": "Sam"
}
`
}
