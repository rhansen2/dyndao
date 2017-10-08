// Package tests is a set of tests that put the various components together and
// demonstrate how they can be combined. (As well as serving as a bit of a test suite...)
package sqlitegen

import (
	"context"
	// These tests are specific to SQLite
	_ "github.com/mattn/go-sqlite3"
	"os"

	"database/sql"
	"fmt"
	"testing"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/orm"
	"github.com/rbastic/dyndao/schema"
)

const PeopleObjectType string = "people"
const AddressesObjectType string = "addresses"

func getDB() *sql.DB {
	sqliteDSN := os.Getenv("SQLITE_DSN")
	if sqliteDSN == "" {
		sqliteDSN = "file::memory:?mode=memory&cache=shared"
	}
	db, err := sql.Open("sqlite3", sqliteDSN)
	if err != nil {
		panic(err)
	}
	return db
}

func TestSaveBasicObject(t *testing.T) {
	var rowsAff int64

	// Configure basic schema
	sch := schema.MockBasicSchema()

	// Retrieve database connection
	db := getDB()
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Initialize ORM
	sqliteORM := orm.New(New("test", sch), sch, db)

	// NOTE: This should force insert
	obj := object.New(PeopleObjectType)
	obj.Set("Name", "Ryan")

	// Create requisite tables
	err := createTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}

	{
		// Save object
		rowsAff, err = sqliteORM.SaveObject(context.TODO(), nil, obj)
		if err != nil {
			t.Fatal(err)
		}
		// Check saved status
		if !obj.GetSaved() {
			t.Fatal("Unknown object error, object not saved")
		}
		if rowsAff == 0 {
			t.Fatal("Rows affected shouldn't be zero initially")
		}
	}

	// Verify personID value
	personID := obj.Get("PersonID").(int64)
	if personID > 1 {
		t.Fatalf("PersonID has the wrong value, has value %d", personID)
	}

	// Test second save to ensure that we don't save the object twice needlessly...
	// This caught a silly bug early on.
	{
		rowsAff, err = sqliteORM.SaveObject(context.TODO(), nil, obj)
		if err != nil {
			t.Fatal(err)
		}
		if rowsAff > 0 {
			t.Fatal("rowsAff should be zero the second time")
		}
	}

	// Now this should force an update, due to name change
	obj.Set("Name", "Joe")

	rowsAff, err = sqliteORM.SaveObject(context.TODO(), nil, obj)
	if err != nil {
		t.Fatal(err)
	}
	// Check saved status
	if rowsAff == 0 {
		t.Fatalf("rowsAff should not be zero")
	}
	if !obj.GetSaved() {
		t.Fatal("system is claiming object isn't saved")
	}

	{
		// Refleshen our object
		latestJoe, err := sqliteORM.Retrieve(context.TODO(), PeopleObjectType, obj.KV)
		if err != nil {
			t.Fatal("retrieve failed: " + err.Error())
		}
		if latestJoe == nil {
			t.Fatal("LatestJoe Should not be nil!")
		}
		if latestJoe.Get("PersonID") != 1 && latestJoe.Get("Name") != "Joe" {
			t.Fatal("latestJoe does not match expectations")
		}
	}

	// End of tests - drop all tables ...
	err = dropTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}
}

func sampleAddressObject() *object.Object {
	addr := object.New(AddressesObjectType)
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
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	sqliteORM := orm.New(New("test", sch), sch, db)
	rootTable := PeopleObjectType

	obj := object.New(rootTable)
	obj.Set("Name", "Ryan")

	addrObj := sampleAddressObject()
	obj.Children[AddressesObjectType] = object.NewArray(addrObj)

	err := createTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: should we rename Save to SaveWithChildren()? anyway, do a complex nested save
	{
		rowsAff, err := sqliteORM.SaveAll(context.TODO(), obj)
		if err != nil {
			t.Fatal("Save:" + err.Error())
		}
		if !obj.GetSaved() {
			t.Fatal("Unknown object error, object not saved")
		}
		if rowsAff == 0 {
			t.Fatal("Rows affected shouldn't be zero initially")
		}
		// Check that children were saved
		for _, childArray := range obj.Children {
			for _, child := range childArray {
				if !child.GetSaved() {
					t.Fatal("Child wasn't saved, type was ", child.Type)
				}
			}
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
			t.Fatal("latestRyan obj does not match expectations")
		}
		// TODO: Verify that addresses are present
	}

	testRetrieveMany(&sqliteORM, t, rootTable)    // test multiple retrieve
	testFleshenChildren(&sqliteORM, t, rootTable) // try fleshen children on person id 1
	testGetParentsViaChild(&sqliteORM, t)

	err = dropTables(db, sch)
	if err != nil {
		t.Fatal(err)
	}

}

func testGetParentsViaChild(o *orm.ORM, t *testing.T) {
	queryVals := make(map[string]interface{})
	queryVals["PersonID"] = 1
	childObj, err := o.Retrieve(context.TODO(), "addresses", queryVals)
	if err != nil {
		t.Fatal(err)
	}

	_, err = o.GetParentsViaChild(context.TODO(), childObj)
	if err != nil {
		t.Fatal(err)
	}
	// TODO: do something better here.
}

func testFleshenChildren(o *orm.ORM, t *testing.T, rootTable string) {
	obj, err := o.Retrieve(context.TODO(), rootTable, map[string]interface{}{
		"PersonID": 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if obj == nil {
		t.Fatal("object should not be nil")
	}
	fleshened, err := o.FleshenChildren(context.TODO(), obj)
	if err != nil {
		t.Fatal(err)
	}
	if fleshened.Type != PeopleObjectType {
		t.Fatal("fleshened object has wrong type, expected", AddressesObjectType)
	}
	if fleshened.Children[AddressesObjectType] == nil {
		t.Fatal("expected children of type ", AddressesObjectType)
	}
	if fleshened.Children["addresses"][0].Get("Address1") != "Test" {
		t.Fatal("expected 'Test' for 'Address1'")
	}
}

func createTables(db *sql.DB, sch *schema.Schema) error {
	gen := New("test", sch)

	for k := range sch.Tables {
		fmt.Println("Creating table ", k)
		sql, err := gen.CreateTable(sch, k)
		fmt.Println("sql->", sql)
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
	gen := New("test", sch)

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

func testRetrieveMany(o *orm.ORM, t *testing.T, rootTable string) {
	// insert another object
	nobj := object.New(rootTable)
	nobj.Set("Name", "Joe")
	{
		rowsAff, err := o.SaveAll(context.TODO(), nobj)
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
	all, err := o.RetrieveMany(context.TODO(), rootTable, make(map[string]interface{}))
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Fatal("Should only be 2 rows inserted")
	}
}
