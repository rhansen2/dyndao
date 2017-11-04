package test

import (
	"fmt"
	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/dyndao/schema/test/mock"
	"testing"
)

func TestDefaultSchema(t *testing.T) {
	sch := schema.DefaultSchema()
	fmt.Println(sch)
}

func TestDefaultTable(t *testing.T) {
	tbl := schema.DefaultTable()
	fmt.Println(tbl)
}

func TestDefaultColumn(t *testing.T) {
	fld := schema.DefaultColumn()
	fmt.Println(fld)
}

func TestDefaultChildTable(t *testing.T) {
	chld := schema.DefaultChildTable()
	fmt.Println(chld)
}

func TestSchemaBasic(t *testing.T) {
	_ = mock.BasicSchema()
}
