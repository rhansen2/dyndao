// Package tests is a set of tests that put the various components together and
// demonstrate how they can be combined. (As well as serving as a bit of a test suite...)
package tests

import (
	"context"
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
	table := "people"

	// NOTE: This should force insert
	obj := object.New(table)
	//obj.Set("PersonID", 1)
	obj.Set("Name", "Ryan")

	db := getDB()
	defer db.Close()

	err := createTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}

	{
		rowsAff, err := orm.Save(context.TODO(), db, sch, obj)
		if err != nil {
			t.Fatal(err)
		}
		if !obj.GetSaved() {
			t.Fatal("Unknown object error, object not saved")
		}
		if rowsAff == 0 {
			t.Fatal("Rows affected shouldn't be zero initially")
		}
	}

	//fmt.Println("PersonID=", obj.Get("PersonID"))

	// Test second save to ensure that we don't save the object twice needlessly...
	// This caught a silly bug early on.
	{
		rowsAff, err := orm.Save(context.TODO(), db, sch, obj)
		if err != nil {
			t.Fatal(err)
		}
		if rowsAff > 0 {
			t.Fatal("rowsAff should be zero the second time")
		}
		//		fmt.Println("rowsAff=", rowsAff)
	}

	// Now this should force an update
	obj.Set("Name", "Joe") // name change

	rowsAff, err := orm.Save(context.TODO(), db, sch, obj)
	if err != nil {
		t.Fatal(err)
	}
	if rowsAff == 0 {
		t.Fatalf("rowsAff should not be zero")
	}

	if !obj.GetSaved() {
		t.Fatal("system is claiming object isn't saved")
	}

	{
		// refleshen our object
		latestJoe, err := orm.Retrieve(context.TODO(), db, sch, table, obj.KV)
		if err != nil {
			t.Fatal("retrieve failed: " + err.Error())
		}
		if latestJoe.Get("PersonID") != 1 && latestJoe.Get("Name") != "Joe" {
			t.Fatal("latestJoe does not match expectations")
		}
		fmt.Println(latestJoe)
	}

	err = dropTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: use contexts down here also

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
