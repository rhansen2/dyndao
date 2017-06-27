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
	// TODO: test all database types that we support.
	db, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		panic(err)
	}
	return db
}

func TestSaveBasicObject(t *testing.T) {
	sch := schema.MockBasicSchema()
	db := getDB()
	defer db.Close()
	sqliteORM := orm.New(sqlitegen.New("test", sch), sch, db)

	table := "people"

	// NOTE: This should force insert
	obj := object.New(table)
	//obj.Set("PersonID", 1)
	obj.Set("Name", "Ryan")


	err := createTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}

	{
		rowsAff, err := sqliteORM.SaveObject(context.TODO(), nil, obj)
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
		rowsAff, err := sqliteORM.SaveObject(context.TODO(), nil, obj)
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

	rowsAff, err := sqliteORM.SaveObject(context.TODO(), nil, obj)
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
		latestJoe, err := sqliteORM.RetrieveObject(context.TODO(), table, obj.KV)
		if err != nil {
			t.Fatal("retrieve failed: " + err.Error())
		}
		if latestJoe == nil {
			t.Fatal("LatestJoe Should not be nil!")
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
	db := getDB()
	defer db.Close()
	sqliteORM := orm.New(sqlitegen.New("test", sch), sch, db)
	rootTable := "people"

	obj := object.New(rootTable)
	obj.Set("Name", "Ryan")

	addrObj := sampleAddressObject()
	obj.Children["addresses"] = object.NewArray(addrObj)

	err := createTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: should we rename Save to SaveWithChildren()? anyway, do a complex nested save
	{
		rowsAff, err := sqliteORM.Save(context.TODO(), obj)
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
		nobj := object.New(rootTable)
		nobj.KV["PersonID"] = obj.Get("PersonID")
		latestRyan, err := sqliteORM.RetrieveWithChildren(context.TODO(), rootTable, obj.KV)
		if err != nil {
			t.Fatal("retrieve failed: " + err.Error())
		}
		if latestRyan.Get("PersonID") != 1 && latestRyan.Get("Name") != "Ryan" {
			t.Fatal("latestRyan does not match expectations")
		}

	}

	{
		queryVals := map[string]interface{}{
			"PersonID": 1,
		}
		childTable := "addresses"
		latestRyan, err := sqliteORM.RetrieveParentViaChild(context.TODO(), childTable, queryVals, nil)
		if err != nil {
			t.Fatal("RetrieveParentViaChild failed: " + err.Error())
		}
		if latestRyan.Get("PersonID") != 1 && latestRyan.Get("Name") != "Ryan" {
			t.Fatal("latestRyan does not match expectations")
		}
		if len(latestRyan.Children) == 0 {
			t.Fatal("latestRyan has no children")
		}
		addrObj := latestRyan.Children["addresses"][0]
		if addrObj == nil {
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
		newJSON, err := mapper.ToJSONFromObject(sch, latestRyan, "{}", "", true)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(newJSON)
	}

	// test multiple retrieve

	{
		// insert another object
		nobj := object.New(rootTable)
		nobj.Set("Name", "Joe")
		{
			rowsAff, err := sqliteORM.Save(context.TODO(), nobj)
			if err != nil {
				t.Fatal("Save:" + err.Error())
			}
			if !nobj.GetSaved() {
				t.Fatal("Unknown object error, object not saved")
			}
			if rowsAff == 0 {
				t.Fatal("Rows affected shouldn't be zero initially")
			}
		}

		// try a full table scan
		all, err := sqliteORM.RetrieveObjects(context.TODO(), rootTable, make(map[string]interface{}))
		if err != nil {
			t.Fatal(err)
		}
		if len(all) != 2 {
			t.Fatal("Should only be 2 rows inserted")
		}
		fmt.Println(all[0])
		fmt.Println(all[1])
		//		fmt.Println(all[2])

		// try fleshen children on person id 1

		{
			obj, err = sqliteORM.RetrieveObject(context.TODO(), rootTable, map[string]interface{}{
				"PersonID": 1,
			})
			if err != nil {
				t.Fatal(err)
			}
			if obj == nil {
				t.Fatal("object should not be nil")
			}
			_, err := sqliteORM.FleshenChildren(context.TODO(), rootTable, obj)
			if err != nil {
				t.Fatal(err)
			}
			// TODO: Fix these tests to actually check the values.... To ensure FleshenChildren works.
			//fmt.Println(obj)
			//fmt.Println(obj.Children["addresses"][0])
		}

		{
			queryVals := make(map[string]interface{})
			queryVals["PersonID"] = 1
			childObj, err := sqliteORM.RetrieveObject(context.TODO(), "addresses", queryVals)
			if err != nil {
				t.Fatal(err)
			}

			objs, err := sqliteORM.GetParentsViaChild(context.TODO(), childObj)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(objs)
		}
	}

	err = dropTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}

}

// TODO: use contexts down here also?

func createTables(db *sql.DB, sch *schema.Schema) error {
	gen := sqlitegen.New("test", sch)

	for k := range sch.Tables {
		fmt.Println("Creating table ", k)
		sql, err := gen.CreateTable(sch, k)
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
	gen := sqlitegen.New("test", sch)

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
