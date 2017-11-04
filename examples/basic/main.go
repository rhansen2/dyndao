package main

import (
	"context"
	"os"
	"time"

	sg "github.com/rbastic/dyndao/sqlgen"

	dorm "github.com/rbastic/dyndao/orm"

	"github.com/rbastic/dyndao/schema/test/mock"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rbastic/dyndao/adapters/core"
	sqliteAdapter "github.com/rbastic/dyndao/adapters/sqlite"
)

var (
	defaultDriver = "sqlite3"
	defaultDSN    = "file::memory:?mode=memory&cache=shared"
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
func getSQLGen() *sg.SQLGenerator {
	sqlGen := sqliteAdapter.New(core.New())
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
	orm := dorm.New(getSQLGen(), sch, db)

	// CreateTables will create all tables within a given schema
	{
		ctx, cancel := getDefaultContext()
		err = orm.CreateTables(ctx)
		cancel()
		fatalIf(err)
	}

	// TODO: put code to work with the database here

	// DropTables will create all tables within a given schema
	{
		ctx, cancel := getDefaultContext()
		err = orm.DropTables(ctx)
		cancel()
		fatalIf(err)
	}
}
