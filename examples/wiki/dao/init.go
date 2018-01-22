package dao

import (
	"context"
	coreAdapter "github.com/rbastic/dyndao/adapters/core"
	sqliteAdapter "github.com/rbastic/dyndao/adapters/sqlite"
	"github.com/rbastic/dyndao/orm"
	sg "github.com/rbastic/dyndao/sqlgen"
	"os"
	"time"
)

var (
	dorm *orm.ORM
)

// TODO: Perhaps I should offer these as standard config opts.
var (
	defaultDriver = "sqlite3"
	defaultDSN    = "file::memory:?mode=memory&cache=shared"
)

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
	sqlGen := sqliteAdapter.New(coreAdapter.New())
	sg.PanicIfInvalid(sqlGen)
	return sqlGen
}

func Init() error {
	driver := getDriver()
	dsn := getDSN()
	db, err := orm.GetDB(driver, dsn)
	if err != nil {
		return err
	}

	sch := InvoiceSchema()
	SetActiveSchema(sch)

	dorm = orm.New(getSQLGen(), sch, db)
	return nil
}

func Stop() {
	dorm.RawConn.Close()
}
