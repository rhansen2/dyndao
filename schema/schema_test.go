package schema

import (
	"fmt"
	"testing"
)

func TestDefaultSchema(t *testing.T) {
	sch := DefaultSchema()
	fmt.Println(sch)
}

func TestDefaultTable(t *testing.T) {
	tbl := DefaultTable()
	fmt.Println(tbl)
}

func TestDefaultField(t *testing.T) {
	fld := DefaultField()
	fmt.Println(fld)
}

func TestDefaultChildTable(t *testing.T) {
	chld := DefaultChildTable()
	fmt.Println(chld)
}

func TestSchemaBasic(t *testing.T) {
	_ = MockBasicSchema()
}
