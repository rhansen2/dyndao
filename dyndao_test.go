// Package tests is a set of tests that put the various components together and
// demonstrate how they can be combined. (As well as serving as a bit of a test suite...)
package tests

import (
	// These tests are specific to SQLite
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"fmt"
	"testing"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	//"github.com/rbastic/dyndao/mapper"
	"github.com/rbastic/dyndao/orm"
	"github.com/rbastic/dyndao/sqlgen/sqlitegen"
)

func getDB() *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		panic(err)
	}
	return db
}

func TestSaveBasicObject(t *testing.T) {
	sch := schema.MockBasicSchema()

	// NOTE: This should force insert
	obj := object.New("people")
	//obj.Set("PersonID", 1)
	obj.Set("Name", "Ryan")

	db := getDB()
	defer db.Close()

	err := createTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}

	rowsAff, err := orm.Save(db, sch, obj)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("rowsAff=", rowsAff)

	fmt.Println("PersonID=", obj.Get("PersonID"))

	err = dropTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}

}

func createTables(db *sql.DB, sch *schema.Schema) error {
	gen := sqlitegen.New("sqlite", "test", sch)

	for k := range sch.Tables {
		fmt.Println("Creating table ", k)
		sql, err := gen.CreateTable(k)
		if err != nil {
			return err
		}
		_, err = db.Exec(sql)
		if err != nil {
			return err
		}

	}
	return nil
}

func dropTables(db *sql.DB, sch *schema.Schema) error {
	gen := sqlitegen.New("sqlite", "test", sch)

	for k := range sch.Tables {
		fmt.Println("Dropping table ", k)
		sql := gen.DropTable(k)
		_, err := db.Exec(sql)
		if err != nil {
			return err
		}

	}
	return nil
}

/*
	func TestLoadBasicObject(t * testing.T) {

	}
	func TestSaveNestedObject(t * testing.T) {

	}

	func TestLoadNestedObject(t * testing.T) {

	}
*/
