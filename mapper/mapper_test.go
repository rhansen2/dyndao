package mapper

import (
	"testing"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/tidwall/gjson"
)

const PeopleObjectType string = "people"
const AddressesObjectType string = "addresses"

func TestBasicMapper(t *testing.T) {
	sch := schema.MockBasicSchema()
	obj, err := ToObjectFromJSON(sch, "", PeopleObjectType, getJSONData())
	if err != nil {
		t.Fatal(err)
	}

	if obj.Type != PeopleObjectType {
		t.Fatal("Incorrect object type, expected people, got " + obj.Type)
	}
	if obj.Get("Name") != "Sam" {
		t.Fatalf("Incorrect object value for key Name, expected Sam, got %s ", obj.Get("Name"))
	}
	if obj.Get("PersonID").(float64) != 1 { // TODO: Fix float64 type.
		t.Fatalf("Incorrect object value for key PersonID, expected 1, got %d", obj.Get("PersonID"))
	}
}

func TestBasicMapperToJSON(t *testing.T) {
	sch := schema.MockBasicSchema()
	obj, err := ToObjectFromJSON(sch, "", PeopleObjectType, getJSONData())
	if err != nil {
		t.Fatal(err)
	}
	json, err := ToJSONFromObject(sch, obj, "", "", true) // TODO: Fix ToJSONFromObject API to not expose setRootPath boolean...
	if err != nil {
		t.Fatal(err)
	}

	testData := getJSONData()
	mappedName := gjson.Get(json, "people.Name").String()
	expectedName := gjson.Get(testData, "people.Name").String()

	if mappedName != expectedName {
		t.Fatalf("Name does not match, expected=%s, got=%s", expectedName, mappedName)
	}

	mappedID := gjson.Get(json, "people.PersonID").Int()
	expectedID := gjson.Get(testData, "people.PersonID").Int()
	if mappedID != expectedID {
		t.Fatalf("PersonID does not match, expected=%d, got=%d", expectedID, mappedID)
	}
}

func getJSONData() string {
	return `{"people": {"Name":"Sam","PersonID":1 } }`
}

func getNestedJSONData() string {
	return `{"people": {"Name":"Sam", "PersonID":1, "addresses": {"Address1":"Test","Address2":"Test2","City":"Nowhere","State":"AZ" } } }`
}

func getNestedObj() *object.Object {
	o := object.New(PeopleObjectType)

	o.KV = map[string]interface{}{
		"PersonID": 1,
		"Name":     "Sam",
	}

	addr := object.New(AddressesObjectType)

	addr.Set("Address1", "Test")
	addr.Set("Address2", "Test2")
	addr.Set("City", "Nowhere")
	addr.Set("State", "AZ")

	o.Children[AddressesObjectType] = object.NewArray(addr)
	return o
}

func TestNestedMapper(t *testing.T) {
	sch := schema.MockNestedSchema()
	obj, err := ToObjectFromJSON(sch, "", PeopleObjectType, getNestedJSONData())
	if err != nil {
		t.Fatal(err)
	}
	json, err := ToJSONFromObject(sch, obj, "", "", true)
	if err != nil {
		t.Fatal(err)
	}

	nj := getNestedJSONData()

	for _, k := range []string{"people.Name", "people.PersonID", "people.addresses.Address1"} {
		gen := gjson.Get(json, k).String()
		exp := gjson.Get(nj, k).String()
		if exp != gen {
			t.Fatalf("data does not match: expected %s, got %s", exp, gen)
		}
	}
}
