// Package db2 is a set of tests that put the various components together and
// demonstrate how they can be combined. (As well as serving as a bit of a test suite...)
//
// In other words, we run database tests, use the generator, use the ORM, etc.
// TODO: More complex test schemas.
package db2

import (
	// Load preferred Oracle driver. Mattn's oci8 had race conditions
	// during testing. rana/ora.v4 has crashes when used for extended
	// testing in server workloads.
	"database/sql"

	_ "bitbucket.org/phiggins/db2cli"

	"fmt"
	"os"
	"testing"

	"github.com/rbastic/dyndao/adapters/core"
	"github.com/rbastic/dyndao/adapters/core/test"
	sg "github.com/rbastic/dyndao/sqlgen"
)

// GetDB is a simple wrapper over sql.Open(), the main purpose being
// to abstract the DSN
func GetDB() *sql.DB {
	dsn := os.Getenv("DB2_DSN")
	if dsn == "" {
		panic("DB2_DSN environment variable is not set, cannot initialize database")
	}
	fmt.Println("Opening with dsn: ", dsn)
	db, err := sql.Open("db2-cli", dsn)
	if err != nil {
		panic(err)
	}
	return db
}

func GetSQLGen() *sg.SQLGenerator {
	sqlGen := core.New()
	sqlGen = New(sqlGen)
	sg.PanicIfInvalid(sqlGen)
	return sqlGen
}

func TestMain(t *testing.T) {
	test.Test(t, GetDB, GetSQLGen)
}
