// Package oraclegen is a set of tests that put the various components together and
// demonstrate how they can be combined. (As well as serving as a bit of a test suite...)
//
// In other words, we run database tests, use the generator, use the ORM, etc.
// TODO: More complex test schemas.
package sqlite

import (
	// Load preferred Oracle driver. Mattn's oci8 had race conditions
	// during testing
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"os"
	"testing"

	"github.com/rbastic/dyndao/adapters/core"
	"github.com/rbastic/dyndao/adapters/core/test"
	sg "github.com/rbastic/dyndao/sqlgen"
)

var (
	defaultDSN = "file::memory:?mode=memory&cache=shared"
)

func GetDB() *sql.DB {
	sqliteDSN := os.Getenv("SQLITE_DSN")
	if sqliteDSN == "" {
		sqliteDSN = defaultDSN
	}
	db, err := sql.Open("sqlite3", sqliteDSN)
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
