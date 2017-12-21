package core

import (
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

func GetLock(g *sg.SQLGenerator, sch *schema.Schema, lockStr string) (string, []interface{}, error) {
	panic("core.GetLock() is not implemented. it must be db-specific")
}

func ReleaseLock(g *sg.SQLGenerator, sch *schema.Schema, lockStr string) (string, []interface{}, error) {
	panic("core.ReleaseLock() is not implemented. it must be db-specific")
}
