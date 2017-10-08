package mapper

import ()

const PeopleObjectType string = "people"
const AddressesObjectType string = "addresses"

func getJSONData() string {
	return `{"people": {"Name":"Sam","PersonID":1 } }`
}

func getNestedJSONData() string {
	return `{"people": {"Name":"Sam", "PersonID":1, "addresses": {"Address1":"Test","Address2":"Test2","City":"Nowhere","State":"AZ" } } }`
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
