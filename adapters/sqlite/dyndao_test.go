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
