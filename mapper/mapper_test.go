package mapper

import (
	"testing"

	"github.com/rbastic/dyndao/schema"
	"github.com/tidwall/gjson"
)

const PeopleObjectType string = "people"
const AddressesObjectType string = "addresses"

func TestBasicMapper(t *testing.T) {
	sch := schema.MockBasicSchema()
	objs, err := ToObjectsFromJSON(sch, getJSONData())
	if err != nil {
		t.Fatal(err)
	}
	if objs == nil {
		t.Fatal("objs is nil")
	}

	obj := objs[0]
	if obj == nil {
		t.Fatal("zeroeth object is nil")
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

func getNestedArrayJSONData() string {
	return `
	{"people": 
		{ "Name":"Sam", 
 	          "PersonID":1, 
		  "addresses": [
			{"Address1":"Test","Address2":"Test2","City":"Nowhere","State":"AZ" },
			{"Address1":"Foo","Address2":"Bar","City":"Lincoln","State":"RI" }
		  ] 
		} 
	}`
}

func TestNestedArrayMapper(t *testing.T) {
	sch := schema.MockNestedSchema()
	objs, err := ToObjectsFromJSON(sch, getNestedArrayJSONData())
	if err != nil {
		t.Fatal(err)
	}
	if len(objs) == 0 {
		t.Fatal("object array is empty")
	}
	obj := objs[0]
	if obj.Type != PeopleObjectType {
		t.Fatal("Object doesn't match expected Type")
	}
	if obj.Get("Name") != "Sam" {
		t.Fatal("Name doesn't match expected value")
	}
	if obj.Get("PersonID").(float64) != 1 {
		t.Fatal("PersonID doesn't match expected value")
	}
	if obj.Children[AddressesObjectType] == nil {
		t.Fatal("Why doesn't our object have any addresses?")
	}
	first := obj.Children[AddressesObjectType][0]
	if first.Get("Address1") == nil {
		t.Fatal("Why is the first Address1 empty?")
	}
	if first.Get("Address1") != "Test" {
		t.Fatal("Why isn't Address1 equal to Test?")
	}
}

func getNestedArrayJSONData2() string {
	return `
	{"people": 
		[
			{ "Name":"Sam", 
			"PersonID":1, 
				"addresses": [
				{"Address1":"Test","Address2":"Test2","City":"Nowhere","State":"AZ" },
				{"Address1":"Foo","Address2":"Bar","City":"Lincoln","State":"RI" }
				] 
			},
			{ "Name":"Ryan", 
			"PersonID":2, 
				"addresses": [
				{"Address1":"Quux","Address2":"TimTams","City":"Nowhere","State":"AZ" },
				{"Address1":"BackAlley","Address2":"Road","City":"Lincoln","State":"RI" }
				] 
			}
		] 
	}`
}

func TestNestedArrayMapper2(t *testing.T) {
	sch := schema.MockNestedSchema()
	objs, err := ToObjectsFromJSON(sch, getNestedArrayJSONData2())
	if err != nil {
		t.Fatal(err)
	}
	if len(objs) == 0 {
		t.Fatal("object array is empty")
	}

	obj := objs[0]
	if obj == nil {
		t.Fatal("object should not be nil")
	}
	if obj.Type != PeopleObjectType {
		t.Fatal("Object doesn't match expected Type")
	}
	if obj.Get("Name") != "Sam" {
		t.Fatal("Name doesn't match expected value")
	}
	if obj.Get("PersonID").(float64) != 1 {
		t.Fatal("PersonID doesn't match expected value")
	}
	if obj.Children[AddressesObjectType] == nil {
		t.Fatal("Why doesn't our object have any addresses?")
	}
	// TODO: Check addresses...
	first := obj.Children[AddressesObjectType][0]
	if first.Get("Address1") == nil {
		t.Fatal("Why is the first Address1 empty?")
	}
	if first.Get("Address1") != "Test" {
		t.Fatal("Why isn't Address1 equal to Test?")
	}

}
