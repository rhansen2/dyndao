package orm

import (
	"github.com/rbastic/dyndao/sqlgen"
)

type ORM struct {
	SQLGen *sqlgen.Generator
}

func New(sgen *sqlgen.Generator) *ORM {
	o := &ORM{SQLGen: sgen}
	return o
}
