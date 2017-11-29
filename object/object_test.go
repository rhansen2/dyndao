package object

import (
	"fmt"
	"testing"
)

func floatTests(t * testing.T, obj * Object) {
	{
		obj.Set("floatTest", 3.141529)
		floatPi, err := obj.GetFloatAlways("floatTest")
		if err != nil {
			panic(err)
		}
		fmt.Println(floatPi)
	}

	{
		obj.Set("intTest", 3141529)
		intPi, err := obj.GetFloatAlways("intTest")
		if err != nil {
			panic(err)
		}
		fmt.Println(intPi)
	}

	{
		intPi, err := obj.GetIntAlways("intTest")
		if err != nil {
			panic(err)
		}
		fmt.Println(intPi)
	}
}

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

	floatTests(t, obj)
}
