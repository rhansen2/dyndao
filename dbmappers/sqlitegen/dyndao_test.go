// Package oraclegen is a set of tests that put the various components together and
// demonstrate how they can be combined. (As well as serving as a bit of a test suite...)
//
// In other words, we run database tests, use the generator, use the ORM, etc.
// TODO: More complex test schemas.
package sqlitegen

import (
	// Load preferred Oracle driver. Mattn's oci8 had race conditions
	// during testing
	"database/sql"
	_ "github.com/mattn/go-sqlite3"


	"os"
	"testing"

	"github.com/rbastic/dyndao/dbmappers/core"
	sg "github.com/rbastic/dyndao/sqlgen"
)

func GetDB() *sql.DB {
	sqliteDSN := os.Getenv("SQLITE_DSN")
	if sqliteDSN == "" {
		sqliteDSN = "file::memory:?mode=memory&cache=shared"
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

func TestMain(t * testing.T) {
	core.Test(t, GetDB, GetSQLGen)
}
