package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"os"
	"testing"

	"github.com/rbastic/dyndao/adapters/core"
	"github.com/rbastic/dyndao/adapters/core/test"
	sg "github.com/rbastic/dyndao/sqlgen"
)

var (
	// TODO: refactor this so it's available from somewhere else
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

// TODO: refactor this so it is available from somewhere else
// (so that user code doesn't have to replicate this)
func GetSQLGen() *sg.SQLGenerator {
	sqlGen := core.New()
	sqlGen = New(sqlGen)
	sg.PanicIfInvalid(sqlGen)
	return sqlGen
}

func TestMain(t *testing.T) {
	test.Test(t, GetDB, GetSQLGen)
}
