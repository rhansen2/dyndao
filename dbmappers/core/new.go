package core

import (
	sg "github.com/rbastic/dyndao/sqlgen"
)

func New() *sg.SQLGenerator {
	g := new(sg.SQLGenerator)
	g.CreateTable = sg.FnCreateTable(CreateTable)
	g.CoreBindingInsert = sg.FnCoreBindingInsert(CoreBindingInsert)
	g.BindingInsert = sg.FnBindingInsert(BindingInsert)
	return g
}

