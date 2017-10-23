package core

import (
	"context"
	"database/sql"

	"testing"

	"github.com/pkg/errors"

	sg "github.com/rbastic/dyndao/sqlgen"
	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/orm"
	"github.com/rbastic/dyndao/schema"
)

const PeopleObjectType string = "people"
const AddressesObjectType string = "addresses"

// function type for GetDB
type FnGetDB func() *sql.DB

// function type for 'GetSQLGenerator'
type FnGetSG func() *sg.SQLGenerator

var (
	GetDB FnGetDB
	getSQLGen FnGetSG
)

func Test( t * testing.T , getDBFn FnGetDB, getSGFN FnGetSG) {
	// Set our functions locally
	GetDB = getDBFn
	getSQLGen = getSGFN

	// Bootstrap the db, run the test suite, drop tables
	TestCreateTables(t)
	TestSuiteNested(t)
	TestDropTables(t)
}

func TestCreateTables(t *testing.T) {
	sch := schema.MockNestedSchema()
	db := GetDB()
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	o := orm.New(getSQLGen(), sch, db)

	err := createTables(o.RawConn, sch)
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
	db := GetDB()
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	// Setup our ORM
	o := orm.New(getSQLGen(), sch, db)
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

	t.Run("Retrieve", func(t *testing.T) {
		// test retrieving the parent, given a child object
		testRetrieve(&o, t, sch)
	})

	t.Run("RetrieveMany", func(t *testing.T) {
		// test multiple retrieve
		testRetrieveMany(&o, t, PeopleObjectType)
	})

	t.Run("FleshenChildren", func(t *testing.T) {
		// try fleshen children on person id 1
		testFleshenChildren(&o, t, PeopleObjectType)
	})

	t.Run("GetParentsViaChild", func(t *testing.T) {
		// test retrieving multiple parents, given a single child object
		testGetParentsViaChild(&o, t)
	})

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

func testRetrieve(o *orm.ORM, t *testing.T, sch *schema.Schema) {
	queryVals := map[string]interface{}{
		"PersonID": 1,
	}
	// refleshen our object
	latestJoe, err := o.Retrieve(context.TODO(), PeopleObjectType, queryVals)
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
	// TODO: How much further should we verify these objects?

}

func testGetParentsViaChild(o *orm.ORM, t *testing.T) {
	// Configure our database query
	queryVals := make(map[string]interface{})
	queryVals["PersonID"] = 1
	// Retrieve a single child object
	childObj, err := o.Retrieve(context.TODO(), "addresses", queryVals)
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
		t.Fatal("expected Addresses children")
	}
	if fleshened.Children["addresses"][0].Get("Address1") != "Test" {
		t.Fatal("expected 'Test' for 'Address1'")
	}
}

func TestDropTables(t *testing.T) {
	sch := schema.MockNestedSchema()
	db := GetDB()
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	o := orm.New(getSQLGen(), sch, db)

	err := dropTables(o.RawConn, sch)
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
	gen := getSQLGen()

	for k := range sch.Tables {
		sql, err := gen.CreateTable(gen, sch, k)
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
	gen := getSQLGen()

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
