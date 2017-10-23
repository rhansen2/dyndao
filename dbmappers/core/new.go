package core

import (
	sg "github.com/rbastic/dyndao/sqlgen"
)

func New() *sg.SQLGenerator {
	g := new(sg.SQLGenerator)
	return g
}

