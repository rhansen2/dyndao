// This is a sample Oracle->MySQL schema migrator. It demonstrates
// how dyndao can be used to translate a schema across databases.
package main

import (
	"fmt"
	"strings"

	sg "github.com/rbastic/dyndao/sqlgen"

	"github.com/rbastic/dyndao/schema/test/mock"

	"github.com/rbastic/dyndao/adapters/core"
	mysqlAdapter "github.com/rbastic/dyndao/adapters/mysql"
)

// TODO: refactor this so it is available from somewhere else
// (so that user code doesn't have to replicate this)
func getSQLGen() *sg.SQLGenerator {
	sqlGen := mysqlAdapter.New(core.New())
	sg.PanicIfInvalid(sqlGen)
	return sqlGen
}

func main() {
	// Load the basic mock schema - load your Oracle schema instead
	sch := mock.BasicSchema()

	for _, t := range sch.Tables {
		for _, f := range t.Columns {
			if strings.ToUpper(f.DBType) == "VARCHAR2" {
				f.DBType = "VARCHAR"
			}

			if strings.ToUpper(f.DBType) == "NUMBER" {
				f.DBType = "INTEGER"
			}

			if f.DBType == "TIMESTAMP" {
				f.DBType = "DATETIME"
				f.Length = 0
			}

			if f.DBType == "CLOB" {
				f.DBType = "TEXT"
			}
		}
	}

	// No ORM instance is necessary here.
	sqlg := getSQLGen()

	for tableName, _ := range sch.Tables {
		{
			s := sqlg.DropTable(tableName)
			fmt.Println(s + ";")
		}

		{
			s, err := sqlg.CreateTable(sqlg, sch, tableName)
			if err != nil {
				panic(err)
			}

			fmt.Println(s + ";")
		}
	}
}
