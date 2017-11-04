package test

import (
	"encoding/json"
	"fmt"
	"github.com/rbastic/dyndao/schema/test/mock"
	"testing"
)

func TestJSONMarshalUnmarshal(t *testing.T) {
	sch := mock.BasicSchema()
	buf, err := json.Marshal(sch)
	if err != nil {
		t.Fatal(err)
	}
	// TODO: Fix this test.
	fmt.Println("Marshalled buf=", string(buf))

	err = json.Unmarshal(buf, &sch)
	if err != nil {
		t.Fatal(err)
	}
	// TODO: Fix this test.
	fmt.Println("Unmarshalled sch=", sch)
}
