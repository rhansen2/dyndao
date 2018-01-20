package oracle

import (
	"context"
	"fmt"
	"os"
	"testing"

	dyndaoORM "github.com/rbastic/dyndao/orm"
	"github.com/rbastic/dyndao/schema"
	_ "gopkg.in/goracle.v2"
)

func fatalEnv() {
	panic("Try setting DRIVER and DSN.\nExample: your_machine$ DRIVER=mysql DSN=root@password//127.0.0.1:3306 go test -v")
}

func fatalIf(t *testing.T, err error) {
	if err != nil {
		panic(err)
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

	dbname := os.Getenv("OWNER")
	if dbname == "" {
		panic("please supply OWNER as an environment parameter. (typically, it's the database username)")
	}
	sch, err := LoadSchema(context.TODO(), db, dbname)
	fatalIf(t, err)

	err = db.Close()
	fatalIf(t, err)

	// TODO: actually validate that the schema was parsed correctly
	// (and try running schema parser against multiple db's)
	//	fmt.Println(sch)
	fmt.Println(sch)

	err = schema.Validate(sch)
	fatalIf(t, err)

	/*
		// TODO: This is some additional schema 'validation' logic. We'd like to manually inspect
		// the metadata that we loaded.
	*/

	for _, tbl := range sch.Tables {
		fmt.Printf("Table %s has %d columns\n", tbl.Name, len(tbl.EssentialColumns))
		for colName, col := range tbl.Columns {
			isIdentityStr := ""
			if col.IsIdentity {
				isIdentityStr = "(IDENTITY)"
			}
			fmt.Printf("   %s%s: %v\n", colName, isIdentityStr, col)
		}
		fmt.Println()
	}
}
