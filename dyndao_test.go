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

	"github.com/rbastic/dyndao/mapper"
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/orm"
	"github.com/rbastic/dyndao/schema"
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
		rowsAff, err := orm.SaveObject(context.TODO(), db, nil, sch, obj)
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

	personID := obj.Get("PersonID").(int64)
	if personID > 1 {
		t.Fatalf("PersonID has the wrong value, has value %d", personID)
	}

	// Test second save to ensure that we don't save the object twice needlessly...
	// This caught a silly bug early on.
	{
		rowsAff, err := orm.SaveObject(context.TODO(), db, nil, sch, obj)
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

	rowsAff, err := orm.SaveObject(context.TODO(), db, nil, sch, obj)
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
		latestJoe, err := orm.RetrieveObject(context.TODO(), db, sch, table, obj.KV)
		if err != nil {
			t.Fatal("retrieve failed: " + err.Error())
		}
		if latestJoe.Get("PersonID") != 1 && latestJoe.Get("Name") != "Joe" {
			t.Fatal("latestJoe does not match expectations")
		}
		//	fmt.Println(latestJoe)
	}

	err = dropTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}
}

func sampleAddressObject() *object.Object {
	addr := object.New("addresses")
	addr.Set("Address1", "Test")
	addr.Set("Address2", "Test2")
	addr.Set("City", "Nowhere")
	addr.Set("State", "AZ")
	addr.Set("Zip", "02865")
	return addr
}

func TestSaveNestedObject(t *testing.T) {
	sch := schema.MockNestedSchema()
	rootTable := "people"

	obj := object.New(rootTable)
	obj.Set("Name", "Ryan")

	addrObj := sampleAddressObject()
	obj.Children["addresses"] = addrObj

	db := getDB()
	defer db.Close()

	err := createTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: should we rename Save to SaveWithChildren()? anyway, do a complex nested save
	{
		rowsAff, err := orm.Save(context.TODO(), db, sch, obj)
		if err != nil {
			t.Fatal("Save:" + err.Error())
		}
		if !obj.GetSaved() {
			t.Fatal("Unknown object error, object not saved")
		}
		if rowsAff == 0 {
			t.Fatal("Rows affected shouldn't be zero initially")
		}
	}

	// now try to do a nested retrieve
	{
		obj.KV["AddressID"] = 1
		latestRyan, err := orm.RetrieveWithChildren(context.TODO(), db, sch, rootTable, obj.KV)
		if err != nil {
			t.Fatal("retrieve failed: " + err.Error())
		}
		if latestRyan.Get("PersonID") != 1 && latestRyan.Get("Name") != "Ryan" {
			t.Fatal("latestRyan does not match expectations")
		}

	}

	{
		queryVals := map[string]interface{}{
			"PersonID":  1,
			"AddressID": 1,
		}
		childTable := "addresses"
		latestRyan, err := orm.RetrieveParentViaChild(context.TODO(), db, sch, childTable, queryVals, nil)
		if err != nil {
			t.Fatal("retrieveparentsviachild failed: " + err.Error())
		}
		if latestRyan.Get("PersonID") != 1 && latestRyan.Get("Name") != "Ryan" {
			t.Fatal("latestRyan does not match expectations")
		}
		if len(latestRyan.Children) == 0 {
			t.Fatal("latestRyan has no children")
		}
		addrObj, ok := latestRyan.Children["addresses"]
		if !ok {
			t.Fatal("latestRyan lacks an 'addresses' child")
		}
		if addrObj != nil {
			if addrObj.Get("Zip") != "02865" {
				t.Fatal("latestRyan has the wrong zipcode")
			}
			if addrObj.Get("City") != "Nowhere" {
				t.Fatal("latestRyan has the wrong city")
			}
			// TODO: write a better expected comparison
		}

		// TODO: Produce nested structure for JSON.
		newJSON, err := mapper.ToJSONFromObject(sch, latestRyan, "{}", "")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(newJSON)

		//	fmt.Println(latestRyan)
		//	fmt.Println(latestRyan.Children["addresses"])
	}

	err = dropTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}

}

// TODO: use contexts down here also?

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

// TODO: Test UPDATE code with an instance of an address
// that now belongs to a different person in the system...
// This will verify that the primary/foreign key code is
// working properly (AKA MultiKey)

/*
	func TestLoadBasicObject(t * testing.T) {

	}
	func TestSaveNestedObject(t * testing.T) {

	}

	func TestLoadNestedObject(t * testing.T) {

	}
*/
