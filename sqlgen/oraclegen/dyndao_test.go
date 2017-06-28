// Package oraclegen is a set of tests that put the various components together and
// demonstrate how they can be combined. (As well as serving as a bit of a test suite...)
package oraclegen

import (
	"context"
	// Load preferred Oracle driver. Mattn's oci8 had race conditions
	// during testing
	_ "gopkg.in/rana/ora.v4"

	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"

	"github.com/rbastic/dyndao/mapper"
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/orm"
	"github.com/rbastic/dyndao/schema"
)

const PeopleObjectType string = "people"
const AddressesObjectType string = "addresses"

// GetDB is a simple wrapper over sql.Open(), the main purpose being
// to abstract the DSN
func GetDB() (*sql.DB, error) {
	// TODO: externalize the DSN and store it in vault
	dsn := os.Getenv("DSN")
	if dsn == "" {
		return nil, errors.New("oracle DSN is not set, cannot initialize database")
	}
	db, err := sql.Open("ora", dsn)
	if err != nil {
		return nil, err
	}
	return db, err
}

func TestCreateTables(t *testing.T) {
	sch := schema.MockNestedSchema()
	db, err := GetDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlGen := New("test", sch)
	o := orm.New(sqlGen, sch, db)

	err = createTables(o.RawConn, sch)
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

func TestSaveBasicObject(t *testing.T) {
	sch := schema.MockNestedSchema()
	db, err := GetDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	o := orm.New(New("test", sch), sch, db)

	table := PeopleObjectType

	// NOTE: This should force insert
	obj := object.New(table)
	//obj.Set("PersonID", 1)
	obj.Set("Name", "Ryan")

	addrObj := sampleAddressObject()
	obj.Children["addresses"] = object.NewArray(addrObj)

	{
		fmt.Println("Saving Ryan")
		rowsAff, err := o.Save(context.TODO(), obj)
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
	if personID != 1 {
		if personID == 2 {
			t.Fatal("Tests are not in a ready state. Pre-existing data is present.")
		}
		t.Fatalf("PersonID has the wrong value, has value %d", personID)
	}

	// Test second save to ensure that we don't save the object twice needlessly...
	// This caught a silly bug early on.
	{
		rowsAff, err := o.SaveObject(context.TODO(), nil, obj)
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

	rowsAff, err := o.SaveObject(context.TODO(), nil, obj)
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
		queryVals := map[string]interface{}{
			"PersonID": 1,
		}
		// refleshen our object
		latestJoe, err := o.RetrieveObject(context.TODO(), table, queryVals)
		if err != nil {
			t.Fatal("retrieve failed: " + err.Error())
		}
		if latestJoe == nil {
			t.Fatal("LatestJoe Should not be nil!")
		}
		fmt.Println(latestJoe)
		if latestJoe.Get("PersonID").(int64) != 1 || latestJoe.Get("Name") != "Joe" {
			t.Fatal("latestJoe does not match expectations")
		}
	}
	testRetrieveParentViaChild(&o, t, sch)

	// test multiple retrieve
	testRetrieveObjects(&o, t, PeopleObjectType)

	// try fleshen children on person id 1
	testFleshenChildren(&o, t, PeopleObjectType)

	testGetParentsViaChild(&o, t)

}

func testRetrieveParentViaChild(o *orm.ORM, t *testing.T, sch *schema.Schema) {
	queryVals := map[string]interface{}{
		"PersonID": 1,
	}
	childTable := "addresses"
	latestRyan, err := o.RetrieveParentViaChild(context.TODO(), childTable, queryVals, nil)
	if err != nil {
		t.Fatal("RetrieveParentViaChild failed: " + err.Error())
	}
	if latestRyan.Get("PersonID").(int64) != 1 && latestRyan.Get("Name") != "Ryan" {
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
	//newJSON, err := mapper.ToJSONFromObject(sch, latestRyan, "{}", "", true)
	_, err = mapper.ToJSONFromObject(sch, latestRyan, "{}", "", true)
	if err != nil {
		t.Fatal(err)
	}
	// TODO: Fix this test.
	//	fmt.Println(newJSON)
}

func testRetrieveObjects(o *orm.ORM, t *testing.T, rootTable string) {
	// insert another object
	nobj := object.New(rootTable)
	nobj.Set("Name", "Joe")
	{
		rowsAff, err := o.Save(context.TODO(), nobj)
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
	all, err := o.RetrieveObjects(context.TODO(), rootTable, make(map[string]interface{}))
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Fatal("Should only be 2 rows inserted")
	}
	/*fmt.Println(all[0])
	fmt.Println(all[1])*/
	//		fmt.Println(all[2])

}

func testGetParentsViaChild(o *orm.ORM, t *testing.T) {
	queryVals := make(map[string]interface{})
	queryVals["PersonID"] = 1
	childObj, err := o.RetrieveObject(context.TODO(), "addresses", queryVals)
	if err != nil {
		t.Fatal(err)
	}

	//objs, err := o.GetParentsViaChild(context.TODO(), childObj)
	_, err = o.GetParentsViaChild(context.TODO(), childObj)
	if err != nil {
		t.Fatal(err)
	}
	// TODO: fix this
	//fmt.Println(objs)
}

func testFleshenChildren(o *orm.ORM, t *testing.T, rootTable string) {
	obj, err := o.RetrieveObject(context.TODO(), rootTable, map[string]interface{}{
		"PersonID": 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if obj == nil {
		t.Fatal("object should not be nil")
	}
	fleshened, err := o.FleshenChildren(context.TODO(), rootTable, obj)
	if err != nil {
		t.Fatal(err)
	}
	if fleshened.Type != PeopleObjectType {
		t.Fatal("fleshened object has wrong type, expected", AddressesObjectType)
	}
	if fleshened.Children[AddressesObjectType] == nil {
		t.Fatal("expected Addresses children")
	}
	if fleshened.Children["addresses"][0].Get("Address1") != "Test" {
		t.Fatal("expected 'Test' for 'Address1'")
	}
	// TODO: Fix these tests to actually check the values.... To ensure FleshenChildren works.
	//fmt.Println(obj)
	//fmt.Println(obj.Children["addresses"][0])
}

func TestDropTables(t *testing.T) {
	sch := schema.MockNestedSchema()
	db, err := GetDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlGen := New("test", sch)
	o := orm.New(sqlGen, sch, db)

	err = dropTables(o.RawConn, sch)
	if err != nil {
		t.Fatal(err)
	}
}

func prepareAndExecSQL(db *sql.DB, sqlStr string) (sql.Result, error) {
	stmt, err := db.PrepareContext(context.TODO(), sqlStr)
	defer stmt.Close()
	r, err := stmt.ExecContext(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "prepareAndExecSQL")
	}
	return r, nil
}

func createTables(db *sql.DB, sch *schema.Schema) error {
	gen := New("test", sch)

	for k := range sch.Tables {
		fmt.Println("Creating table ", k)
		sql, err := gen.CreateTable(sch, k)
		fmt.Println("CreateTable SQL", sql)
		if err != nil {
			return err
		}
		r, err := prepareAndExecSQL(db, sql)
		if err != nil {
			return errors.Wrap(err, "createTables")
		}
		rowsAff, err := r.RowsAffected()
		if err != nil {
			return err
		}
		fmt.Println("RowsAffected=", rowsAff)

	}
	return nil
}

func dropTables(db *sql.DB, sch *schema.Schema) error {
	gen := New("test", sch)

	for k := range sch.Tables {
		fmt.Println("Dropping table ", k)
		sql := gen.DropTable(k)
		fmt.Println("DropTable SQL", sql)
		r, err := prepareAndExecSQL(db, sql)
		if err != nil {
			return errors.Wrap(err, "dropTables")
		}
		rowsAff, err := r.RowsAffected()
		if err != nil {
			return err
		}
		fmt.Println("RowsAffected=", rowsAff)

	}
	return nil
}
