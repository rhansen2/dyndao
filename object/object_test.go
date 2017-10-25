package object

import (
	"fmt"
	"testing"
)

func TestObject(t *testing.T) {
	obj := New("person")

	obj.Set("name", "Ryan")
	obj.Set("age", 30)
	obj.Set("date_of_birth", "01-01-1970")
	obj.Set("favorite_number", 3.141529)

	fmt.Println("obj is ", obj)

	obj.ResetChangedColumns()

	obj.Set("age", 31)
	obj.Set("id", NewSQLValue("SYS_GUID()"))

	fmt.Println(obj.ChangedColumns)
	fmt.Println(obj.KV)
}
