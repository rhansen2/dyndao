package infoschema

import (
	"context"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	dyndaoORM "github.com/rbastic/dyndao/orm"
	"github.com/rbastic/dyndao/schema"
	//"database/sql"
)

func fatalEnv() {
	panic("Try setting DRIVER and DSN.\nExample: your_machine$ DRIVER=mysql DSN=root@password//127.0.0.1:3306 go test -v")
}

func fatalIf(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestBasicIS(t *testing.T) {
	driver := os.Getenv("DRIVER")
	dsn := os.Getenv("DSN")
	if driver == "" || dsn == "" {
		fatalEnv()
	}

	db, err := dyndaoORM.GetDB(driver, dsn)
	fatalIf(t, err)

	// TODO: use a different schema than 'test'
	sch, err := LoadSchema(context.TODO(), db, "test")
	fatalIf(t, err)

	err = db.Close()
	fatalIf(t, err)

	// TODO: actually validate that the schema was parsed correctly
	// (and try running schema parser against multiple db's)
	//	fmt.Println(sch)

	err = schema.Validate(sch)
	fatalIf(t, err)
}
