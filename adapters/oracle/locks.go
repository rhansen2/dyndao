package oracle

import (
	"fmt"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

func GetLock(g *sg.SQLGenerator, sch *schema.Schema, lockStr string) (string, []interface{}, error) {
	return fmt.Sprintf("BEGIN request_lock('%s'); END;", lockStr), nil, nil
}

func ReleaseLock(g *sg.SQLGenerator, sch *schema.Schema, lockStr string) (string, []interface{}, error) {
	return fmt.Sprintf("BEGIN release_lock('%s'); END;", lockStr), nil, nil
}
