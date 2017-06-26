package mapper

import (
	"fmt"
	"testing"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

func TestBasicMapper(t *testing.T) {
	sch := schema.MockBasicSchema()

	obj, err := ToObjectFromJSON(sch, "people", getJSONData())

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(obj)
}

func TestBasicMapperToJSON(t *testing.T) {
	sch := schema.MockBasicSchema()

	obj, err := ToObjectFromJSON(sch, "people", getJSONData())

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(obj)

	json, err := ToJSONFromObject(sch, obj, "", "")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(getJSONData())
	fmt.Println(json)
	if json != getJSONData() {
		t.Fatal("json data doesn't match")
	}
}

func getJSONData() string {
	return `{"Name":"Sam","PersonID":1}`
}

// TODO: Woops, this is not actually a 'nested object structure'...
// This would be 'flat json to nested object'...
func getNestedJSONData() string {
	return `{"Name":"Sam","PersonID":1,"Address1":"Test","Address2":"Test2","City":"Nowhere","State":"AZ" }`
}

func getNestedObj() *object.Object {
	o := object.New("people")

	// TODO: use constants for object Types instead of hard-coded strings
	o.KV = map[string]interface{}{
		"PersonID": 1,
		"Name":     "Sam",
	}

	addr := object.New("addresses")

	// TODO: use constants for object Types instead of hard-coded strings
	addr.Set("Address1", "Test")
	addr.Set("Address2", "Test2")
	addr.Set("City", "Nowhere")
	addr.Set("State", "AZ")

	o.Children["addresses"] = object.NewObjectArray(addr)
	return o
}

func TestNestedMapper(t *testing.T) {
	sch := schema.MockNestedSchema()
	obj, err := ToObjectFromJSON(sch, "people", getNestedJSONData())

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(obj)
	fmt.Println(obj.Children)
	fmt.Println(obj.Children["addresses"])

	testObj := getNestedObj()
	fmt.Println(testObj)
	fmt.Println(testObj.Children)
	fmt.Println(testObj.Children["addresses"])

	json, err := ToJSONFromObject(sch, obj, "", "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(json)
}
