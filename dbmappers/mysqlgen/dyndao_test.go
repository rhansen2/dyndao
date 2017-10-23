package mysqlgen

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"os"
	"testing"

	"github.com/rbastic/dyndao/dbmappers/core"
	sg "github.com/rbastic/dyndao/sqlgen"
)

// GetDB is a simple wrapper over sql.Open(), the main purpose being
// to abstract the DSN

func GetDB() *sql.DB {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		panic("MYSQL_DSN environment variable is not set, cannot initialize database")
	}
	db, err := sql.Open("mysql", dsn)
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
	core.Test(t, GetDB, GetSQLGen)
}
