// Package oraclegen is a set of tests that put the various components together and
// demonstrate how they can be combined. (As well as serving as a bit of a test suite...)
//
// In other words, we run database tests, use the generator, use the ORM, etc.
// TODO: More complex test schemas.
package oraclegen

import (
	"context"
	// Load preferred Oracle driver. Mattn's oci8 had race conditions
	// during testing
	"database/sql"
	_ "gopkg.in/goracle.v2"

	"os"
	"testing"

	"github.com/pkg/errors"

	"github.com/rbastic/dyndao/mapper"
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/orm"
	"github.com/rbastic/dyndao/schema"
	"github.com/tidwall/gjson"
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
	db, err := sql.Open("goracle", dsn)
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
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	sqlGen := New("test", sch, false)
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
func makeDefaultPersonWithAddress() *object.Object {
	// NOTE: This should force insert
	obj := object.New(PeopleObjectType)
	//obj.Set("PersonID", 1)
	obj.Set("Name", "Ryan")

	addrObj := sampleAddressObject()
	obj.Children["addresses"] = object.NewArray(addrObj)
	return obj
}

func TestSuiteNested(t *testing.T) {
	// Test schema
	sch := schema.MockNestedSchema()
	// Grab database connection
	db, err := GetDB()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	// Setup our ORM
	o := orm.New(New("test", sch, false), sch, db)
	// Construct our default mock object
	obj := makeDefaultPersonWithAddress()
	// Save our default object
	t.Run("SaveMockObject", func(t *testing.T) {
		saveMockObject(t, &o, obj)
	})
	// Validate that we correctly fleshened the primary key
	t.Run("ValidatePersonID", func(t *testing.T) {
		validatePersonID(t, obj)
		// TODO: Make sure we saved the Address with a person id also
	})

	// Validate that we correctly saved the children
	t.Run("ValidateChildrenSaved", func(t *testing.T) {
		validateChildrenSaved(t, obj)
	})

	// Test second additional Save to ensure that we don't save
	// the object twice needlessly... This caught a silly bug early on.
	t.Run("TestAdditionalSave", func(t *testing.T) {
		rowsAff, err := o.SaveObject(context.TODO(), nil, obj)
		if err != nil {
			t.Fatal(err)
		}
		if rowsAff != 0 {
			t.Fatal("rowsAff should be zero the second time")
		}
	})

	// Now, trigger an update.
	t.Run("TestUpdateObject", func(t *testing.T) {
		// Changing the name should become an
		// update when we save
		obj.Set("Name", "Joe")
		// Test saving the object
		testSaveObject(&o, t, obj)
	})

	t.Run("RetrieveObject", func(t *testing.T) {
		// test retrieving the parent, given a child object
		testRetrieveObject(&o, t, sch)
	})

	t.Run("RetrieveObjects", func(t *testing.T) {
		// test multiple retrieve
		testRetrieveObjects(&o, t, PeopleObjectType)
	})

	t.Run("FleshenChildren", func(t *testing.T) {
		// try fleshen children on person id 1
		testFleshenChildren(&o, t, PeopleObjectType)
	})

	t.Run("GetParentsViaChild", func(t *testing.T) {
		// test retrieving multiple parents, given a single child object
		testGetParentsViaChild(&o, t)
	})

	// JSON mapper tests
	var newJSON string
	t.Run("JSONMapper", func(t *testing.T) {
		latestRyan, err := o.RetrieveObject(context.TODO(), PeopleObjectType,
			map[string]interface{}{
				"PersonID": 1,
			})
		if err != nil {
			t.Fatal(err)
		}
		newJSON, err = mapper.ToJSONFromObject(sch, latestRyan, "{}", "", true)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("JSONMapperValidate", func(t *testing.T) {
		validateJSONMapper(t, newJSON)
	})

	t.Run("JSONMapperFrom", func(t *testing.T) {
		objs, err := mapper.ToObjectsFromJSON(sch, newJSON)
		if err != nil {
			t.Fatal(err)
		}
		if objs == nil {
			t.Fatal("objs was nil")
		}
		if len(objs) != 1 {
			t.Fatal("objs length was: ", len(objs), "expected 1")
		}
		obj := objs[0]
		if obj == nil {
			t.Fatal("obj was nil")
		}
		validateJSONMapperFrom(t, obj)
	})

	// TODO: More JSON mapper tests <-> (both To and From)
}

func validateJSONMapper(t *testing.T, json string) {
	if gjson.Get(json, "people.Name").String() != "Joe" {
		t.Fatal("people.Name was not Joe")
	}
	if gjson.Get(json, "people.PersonID").String() != "1" {
		t.Fatal("people.PersonID was not 1")
	}
}

func validateJSONMapperFrom(t *testing.T, obj *object.Object) {
	if obj.Type != PeopleObjectType {
		t.Fatal("Object is wrong type")
	}
	pn := obj.Get("Name")
	if pn != "Joe" {
		t.Fatal("Object has wrong name")
	}

	pi := obj.Get("PersonID")
	if pi.(float64) != 1 {
		t.Fatal("Object has wrong PersonID")
	}
}

func saveMockObject(t *testing.T, o *orm.ORM, obj *object.Object) {
	rowsAff, err := o.SaveAll(context.TODO(), obj)
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

func validatePersonID(t *testing.T, obj *object.Object) {
	personID, err := obj.GetIntAlways("PersonID")
	if err != nil {
		t.Fatalf("Had problems retrieving PersonID as int: %s", err.Error())
	}
	if personID != 1 {
		if personID == 2 {
			t.Fatal("Tests are not in a ready state. Pre-existing data is present.")
		}
		t.Fatalf("PersonID has the wrong value, has value %d", personID)
	}
}

func validateChildrenSaved(t *testing.T, obj *object.Object) {
	for _, childs := range obj.Children {
		for _, v := range childs {
			if !v.GetSaved() {
				t.Fatal("Child object wasn't saved")
			}
			addrID := v.Get("AddressID").(int64)
			if addrID < 1 {
				t.Fatal("AddressID was not what we expected")
			}
		}
	}
}

func testSaveObject(o *orm.ORM, t *testing.T, obj *object.Object) {
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

}

func testRetrieveObject(o *orm.ORM, t *testing.T, sch *schema.Schema) {
	queryVals := map[string]interface{}{
		"PersonID": 1,
	}
	// refleshen our object
	latestJoe, err := o.RetrieveObject(context.TODO(), PeopleObjectType, queryVals)
	if err != nil {
		t.Fatal("retrieve failed: " + err.Error())
	}
	if latestJoe == nil {
		t.Fatal("LatestJoe Should not be nil!")
	}
	// TODO: Do a common refactor on this sort of code
	if latestJoe.Get("PersonID").(int64) != 1 || latestJoe.Get("Name") != "Joe" {
		t.Fatal("latestJoe does not match expectations")
	}
}

func testRetrieveObjects(o *orm.ORM, t *testing.T, rootTable string) {
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
	all, err := o.RetrieveObjects(context.TODO(), rootTable, make(map[string]interface{}))
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Fatal("Should only be 2 rows inserted")
	}
	// TODO: How much further should we verify these objects?

}

func testGetParentsViaChild(o *orm.ORM, t *testing.T) {
	// Configure our database query
	queryVals := make(map[string]interface{})
	queryVals["PersonID"] = 1
	// Retrieve a single child object
	childObj, err := o.RetrieveObject(context.TODO(), "addresses", queryVals)
	if err != nil {
		t.Fatal(err)
	}
	// A silly test, but verifies basic assumptions
	if childObj.Type != "addresses" {
		t.Fatal("Unknown child object retrieved", childObj)
	}
	// Retrieve the parents of that child object
	objs, err := o.GetParentsViaChild(context.TODO(), childObj)
	if err != nil {
		t.Fatal(err)
	}
	// Validate expected data
	if len(objs) != 1 {
		t.Fatal("Unknown length of objs, expected 1, got ", len(objs))
	}
	obj := objs[0]
	if obj.Get("PersonID").(int64) != 1 || obj.Get("Name").(string) != "Joe" {
		t.Fatal("Object values differed from expectations", obj)
	}
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
	fleshened, err := o.FleshenChildren(context.TODO(), obj)
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
}

func TestDropTables(t *testing.T) {
	sch := schema.MockNestedSchema()
	db, err := GetDB()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	sqlGen := New("test", sch, false)
	o := orm.New(sqlGen, sch, db)

	err = dropTables(o.RawConn, sch)
	if err != nil {
		t.Fatal(err)
	}
}

func prepareAndExecSQL(db *sql.DB, sqlStr string) (sql.Result, error) {
	stmt, err := db.PrepareContext(context.TODO(), sqlStr)
	if err != nil {
		return nil, errors.Wrap(err, "prepareAndExecSQL/PrepareContext")
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}()
	r, err := stmt.ExecContext(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "prepareAndExecSQL/ExecContext")
	}
	return r, nil
}

func createTables(db *sql.DB, sch *schema.Schema) error {
	gen := New("test", sch, false)

	for k := range sch.Tables {
		sql, err := gen.CreateTable(sch, k)
		if err != nil {
			return err
		}
		_, err = prepareAndExecSQL(db, sql)
		if err != nil {
			return errors.Wrap(err, "createTables")
		}
	}
	return nil
}

func dropTables(db *sql.DB, sch *schema.Schema) error {
	gen := New("test", sch, false)

	for k := range sch.Tables {
		sql := gen.DropTable(k)
		r, err := prepareAndExecSQL(db, sql)
		if err != nil {
			return errors.Wrap(err, "dropTables")
		}
		_, err = r.RowsAffected()
		if err != nil {
			return err
		}
	}
	return nil
}
