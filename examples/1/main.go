package main

import (
	"context"
	"os"

	sg "github.com/rbastic/dyndao/sqlgen"

	"github.com/rbastic/dyndao/object"
	dorm "github.com/rbastic/dyndao/orm"

	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/dyndao/schema/mock"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rbastic/dyndao/adapters/core"
	"github.com/rbastic/dyndao/adapters/sqlite"
)

var (
	defaultDriver = "sqlite3"
	defaultDSN = "file::memory:?mode=memory&cache=shared"
)

func fatalIf(err error) {
	if err != nil {
		panic(err)
	}
}

func getDefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func getDriver() string {
	return defaultDriver
}

func getDSN() string {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = defaultDSN
	}

	return dsn
}

// TODO: refactor this so it is available from somewhere else
// (so that user code doesn't have to replicate this)
func GetSQLGen() *sg.SQLGenerator {
	sqlGen := core.New()
	sqlGen = New(sqlGen)
	sg.PanicIfInvalid(sqlGen)
	return sqlGen
}

func main() {
	driver := getDriver()
	dsn := getDSN()
	db, err := dorm.GetDB(driver, dsn)
	fatalIf(err)

	if db == nil {
		panic("empty database connection received")
	}

	defer func() {
		err = db.Close()
		fatalIf(err)
	}()

	sch := mock.NestedSchema()
	orm := dorm.New(GetSQLGen(), sch, db)

	
}
