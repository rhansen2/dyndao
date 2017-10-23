package core

import (
	sg "github.com/rbastic/dyndao/sqlgen"
)

func New() *sg.SQLGenerator {
	g := new(sg.SQLGenerator)
	g.CreateTable = sg.FnCreateTable(CreateTable)
	return g
}

