// Package test is the core test suite for all database drivers.
// A conformant implementation should pass all tests.
//
// This test suite is not meant to be run on it's own. The individual
// driver folders have their own dyndao_test.go which bootstraps this
// code.

/*
	TODO: I need to research Go testing patterns more thoroughly.
	For database driver testing, it would appear that panic()
	would be more suitable than t.Fatal() or t.Fatalf()...

	Then I would get a stacktrace, and program execution would halt,
	making it much more apparent that something had gone wrong.
*/

package test

import (
	"context"
	"database/sql"
	"os"
	"reflect"
	"time"
	//"fmt"

	"sync"
	"testing"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/orm"
	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/dyndao/schema/test/mock"

	sg "github.com/rbastic/dyndao/sqlgen"
)

type FnGetDB func() *sql.DB          // function type for GetDB
type FnGetSG func() *sg.SQLGenerator // function type for GetSQLGenerator

const (
	ColPersonID = "PersonID"
)

var (
	globalPersonID int64
)

var (
	GetDB     FnGetDB
	getSQLGen FnGetSG

	// Because CreateTables now has the potential to side-effect the
	// schema, modifying database types based on dyndao's perceived type
	// affinity, we want to keep a stable copy of that laying around.
	// -TODO- Implement / support per-adapter type mappers?
	cachedSchema *schema.Schema
)

func getSchema() *schema.Schema {
	if cachedSchema == nil {
		cachedSchema = mock.NestedSchema()
	}
	return cachedSchema
}

func getDefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 2*time.Second)
}

// Longer context for race condition testing
func getRaceContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func fatalIf(err error) {
	if err != nil {
		panic(err)
	}
}

func PingCheck(t *testing.T, db *sql.DB) {
	ctx, cancel := getDefaultContext()
	err := db.PingContext(ctx)
	cancel()
	fatalIf(err)
}

func dirtyTest(obj *object.Object) {
	if obj.IsDirty() {
		panic("system claims object was not saved")
	}
}

func Test(t *testing.T, getDBFn FnGetDB, getSGFN FnGetSG) {
	// Set our functions locally
	GetDB = getDBFn
	getSQLGen = getSGFN

	db := GetDB()
	if db == nil {
		t.Fatal("dyndao: core/test/Test: GetDB() returned a nil value.")
	}
	defer func() {
		err := db.Close()
		fatalIf(err)
	}()

	// TODO: for testing
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(9) // ???

	// Bootstrap the db, run the test suite, drop tables
	t.Run("TestPingCheck", func(t *testing.T) {
		PingCheck(t, db)
	})

	// If the test suite fails before it gets to drop all tables, you'll
	// want to pass DROP_TABLES=1 before rerunning so that they're wiped.
	if os.Getenv("DROP_TABLES") != "" {
		t.Run("TestDropTables", func(t *testing.T) {
			TestDropTables(t, db)
		})
	}

	t.Run("TestCreateTables", func(t *testing.T) {
		TestCreateTables(t, db)
	})

	TestSuiteNested(t, db)

	t.Run("TestDropTables", func(t *testing.T) {
		TestDropTables(t, db)
	})
}

func TestCreateTables(t *testing.T, db *sql.DB) {
	sch := getSchema()
	o := orm.New(getSQLGen(), sch, db)

	ctx, cancel := getDefaultContext()
	err := o.CreateTables(ctx)
	cancel()
	if err != nil {
		panic(err)
	}
}

func TestDropTables(t *testing.T, db *sql.DB) {
	sch := getSchema()
	o := orm.New(getSQLGen(), sch, db)

	ctx, cancel := getDefaultContext()
	err := o.DropTables(ctx)
	cancel()
	fatalIf(err)
}

func validateMock(t *testing.T, obj *object.Object) {
	// Validate that we correctly fleshened the primary key
	t.Run("ValidatePerson/ID", func(t *testing.T) {
		validatePersonID(t, obj)

		// TODO: Make sure we saved the Address with a person id also
	})
	t.Run("ValidatePerson/NullText", func(t *testing.T) {
		validateNullText(t, obj)
	})

	t.Run("ValidatePerson/NullVarchar", func(t *testing.T) {
		validateNullVarchar(t, obj)
	})

	t.Run("ValidatePerson/NullInt", func(t *testing.T) {
		validateNullInt(t, obj)
	})

	t.Run("ValidatePerson/NullFloat", func(t *testing.T) {
		validateNullFloat(t, obj)
	})
	t.Run("ValidatePerson/NullBlob", func(t *testing.T) {
		validateNullBlob(t, obj)
	})
	t.Run("ValidatePerson/NullTimestamp", func(t *testing.T) {
		validateNullBlob(t, obj)
	})

	// Validate that we correctly saved the children
	t.Run("ValidateChildrenSaved", func(t *testing.T) {
		validateChildrenSaved(t, obj)
	})
}

func TestSuiteNested(t *testing.T, db *sql.DB) {
	sch := getSchema()
	o := orm.New(getSQLGen(), sch, db)     // Setup our ORM
	obj := mock.DefaultPersonWithAddress() // Construct our default mock object

	// Save our default object
	t.Run("SaveMockObject", func(t *testing.T) {
		saveMockObject(o, t, obj)
		var err error
		globalPersonID, err = obj.GetIntAlways(ColPersonID)
		if err != nil {
			t.Fatal("Encountered err when attempting to read PersonID:", err)
		}
	})

	validateMock(t, obj)

	// Once we get basic saving working, do some race condition testing
	if os.Getenv("TEST_RACE") != "" {
		t.Run("RaceConditionSave", func(t *testing.T) {
			testRaceConditionSave(o, t, mock.PeopleObjectType)
		})
	}

	// Test second additional Save to ensure that we don't save
	// the object twice needlessly... This caught a silly bug early on.
	t.Run("TestAdditionalSave", func(t *testing.T) {
		ctx, cancel := getDefaultContext()
		rowsAff, err := o.Save(ctx, nil, obj)
		cancel()
		fatalIf(err)
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
		ctx, cancel := getDefaultContext()
		testSave(ctx, o, t, obj)
		cancel()
	})

	t.Run("Retrieve", func(t *testing.T) {
		// test retrieving the parent, given a child object
		testRetrieve(o, t, sch)
	})

	t.Run("RetrieveMany", func(t *testing.T) {
		// test multiple retrieve
		testRetrieveMany(o, t, mock.PeopleObjectType)
	})

	t.Run("FleshenChildren", func(t *testing.T) {
		// try fleshen children on person id 1
		testFleshenChildren(o, t, mock.PeopleObjectType)
	})

	t.Run("GetParentsViaChild", func(t *testing.T) {
		// test retrieving multiple parents, given a single child object
		testGetParentsViaChild(o, t)
	})

	t.Run("Delete", func(t *testing.T) {
		// test mock object delete
		testDeleteMockObject(o, t, obj)
	})
}

func testDeleteMockObject(o *orm.ORM, t *testing.T, obj *object.Object) {
	ctx, cancel := getDefaultContext()
	rowsAff, err := o.Delete(ctx, nil, obj)
	cancel()
	fatalIf(err)

	dirtyTest(obj)

	if rowsAff == 0 {
		t.Fatal("Rows affected shouldn't be zero")
	}
}

func saveMockObject(o *orm.ORM, t *testing.T, obj *object.Object) {
	ctx, cancel := getDefaultContext()
	rowsAff, err := o.SaveAll(ctx, obj)
	cancel()
	fatalIf(err)

	dirtyTest(obj)

	if rowsAff == 0 {
		t.Fatal("Rows affected shouldn't be zero initially")
	}
}

func validatePersonID(t *testing.T, obj *object.Object) {
	_, err := obj.GetIntAlways(ColPersonID)
	if err != nil {
		t.Fatalf("Had problems retrieving PersonID as int: %s", err.Error())
	}
}

func validateNullText(t *testing.T, obj *object.Object) {
	if !obj.ValueIsNULL(obj.Get("NullText")) {
		t.Fatal("validateNullText expected NULL value")
	}
}

func validateNullInt(t *testing.T, obj *object.Object) {
	if !obj.ValueIsNULL(obj.Get("NullInt")) {
		t.Fatal("validateNullInt: expected NULL value")
	}
}

func validateNullFloat(t *testing.T, obj *object.Object) {
	if !obj.ValueIsNULL(obj.Get("NullFloat")) {
		t.Fatal("validateNullFloat: expected NULL value")
	}
}

func validateNullVarchar(t *testing.T, obj *object.Object) {
	if !obj.ValueIsNULL(obj.Get("NullVarchar")) {
		t.Fatal("validateNullVarchar: expected NULL value")
	}
}

func validateNullBlob(t *testing.T, obj *object.Object) {
	if !obj.ValueIsNULL(obj.Get("NullBlob")) {
		t.Fatal("validateNullBlob: expected NULL value")
	}
}

func validateNullTimestamp(t *testing.T, obj *object.Object) {
	if !obj.ValueIsNULL(obj.Get("NullTimestamp")) {
		t.Fatal("validateNullTimestamp: expected NULL value")
	}
}

func validateChildrenSaved(t *testing.T, obj *object.Object) {
	for _, childs := range obj.Children {
		for _, v := range childs {
			dirtyTest(v)

			addrID := v.Get("AddressID").(int64)
			if addrID < 1 {
				t.Fatal("AddressID was not what we expected")
			}
		}
	}
}

func testDelete(ctx context.Context, o *orm.ORM, t *testing.T, obj *object.Object) {
	rowsAff, err := o.Delete(ctx, nil, obj)

	fatalIf(err)
	if rowsAff == 0 {
		t.Fatal("rowsAff should not be zero")
	}
	if rowsAff != 1 {
		t.Fatalf("rowsAff should not be %d, expected 1", rowsAff)
	}

	dirtyTest(obj)
}

func testSave(ctx context.Context, o *orm.ORM, t *testing.T, obj *object.Object) {
	rowsAff, err := o.Save(ctx, nil, obj)

	fatalIf(err)
	if rowsAff == 0 {
		t.Fatalf("rowsAff should not be zero")
	}

	dirtyTest(obj)
}

func testRetrieve(o *orm.ORM, t *testing.T, sch *schema.Schema) {
	queryVals := map[string]interface{}{
		ColPersonID: globalPersonID,
	}
	// refleshen our object
	ctx, cancel := getDefaultContext()
	latestJoe, err := o.Retrieve(ctx, mock.PeopleObjectType, queryVals)
	cancel()
	if err != nil {
		t.Fatal("retrieve failed: " + err.Error())
	}
	if latestJoe == nil {
		t.Fatal("LatestJoe Should not be nil!")
	}
	// TODO: Do a common refactor on this sort of code
	nameStr, err := latestJoe.GetStringAlways("Name")
	fatalIf(err)
	if latestJoe.Get(ColPersonID).(int64) != globalPersonID || nameStr != "Joe" {
		t.Fatal("latestJoe does not match expectations")
	}
}

func testRetrieveMany(o *orm.ORM, t *testing.T, rootTable string) {
	// TODO: Refactor this?

	// insert another object
	nobj := object.New(rootTable)
	nobj.Set("Name", "Joe")
	{
		ctx, cancel := getDefaultContext()
		rowsAff, err := o.SaveAll(ctx, nobj)
		cancel()
		if err != nil {
			t.Fatal("Save:" + err.Error())
		}
		dirtyTest(nobj)

		if rowsAff == 0 {
			t.Fatal("Rows affected shouldn't be zero initially")
		}
	}

	// try a full table scan
	ctx, cancel := getDefaultContext()
	all, err := o.RetrieveMany(ctx, rootTable, make(map[string]interface{}))
	cancel()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Fatal("Should only be 2 rows inserted")
	}

	car := all[0]
	cdr := all[1]

	{
		if car.Get("Name") == cdr.Get("Name") && car.Get(ColPersonID) != cdr.Get(ColPersonID) {
			// pass
		} else {
			t.Fatal("objects weren't what we expected? they appear to be the same?")
		}
	}

	if reflect.DeepEqual(car, cdr) {
		t.Fatal("Objects matched, this was not expected")
	}
}

func testGetParentsViaChild(o *orm.ORM, t *testing.T) {
	// Configure our database query
	queryVals := make(map[string]interface{})
	queryVals[ColPersonID] = globalPersonID

	// Retrieve a single child object
	ctx, cancel := getDefaultContext()
	childObj, err := o.Retrieve(ctx, "addresses", queryVals)
	cancel()

	fatalIf(err)
	if childObj == nil {
		t.Fatal("testGetParentsViaChild: childObj was nil")
	}
	// A silly test, but verifies basic assumptions
	if childObj.Type != "addresses" {
		t.Fatal("Unknown child object retrieved", childObj)
	}

	// Retrieve the parents of that child object
	ctx, cancel = getDefaultContext()
	objs, err := o.GetParentsViaChild(ctx, childObj)
	fatalIf(err)

	// Validate expected data
	if len(objs) != 1 {
		t.Fatal("Unknown length of objs, expected 1, got ", len(objs))
	}
	obj := objs[0]

	nameStr, err := obj.GetStringAlways("Name")
	fatalIf(err)
	if obj.Get(ColPersonID).(int64) != globalPersonID || nameStr != "Joe" {
		t.Fatal("Object values differed from expectations", obj)
	}
}

func testFleshenChildren(o *orm.ORM, t *testing.T, rootTable string) {
	ctx, cancel := getDefaultContext()
	obj, err := o.Retrieve(ctx, rootTable, map[string]interface{}{
		ColPersonID: globalPersonID,
	})
	cancel()
	fatalIf(err)

	if obj == nil {
		t.Fatal("object should not be nil")
	}

	{
		ctx, cancel := getDefaultContext()
		fleshened, err := o.FleshenChildren(ctx, obj)
		cancel()
		fatalIf(err)

		if fleshened.Type != mock.PeopleObjectType {
			t.Fatal("fleshened object has wrong type, expected", mock.AddressesObjectType)
		}
		if fleshened.Children[mock.AddressesObjectType] == nil {
			t.Fatal("expected Addresses children")
		}
		address1, err := fleshened.Children["addresses"][0].GetStringAlways("Address1")
		if err != nil {
			panic(err)
		}
		expectedStr := "Test"
		if address1 != "Test" {
			t.Fatalf("expected %s for 'Address1', address1 was %s", expectedStr, address1)
		}
	}
}

func testRaceConditionSave(o *orm.ORM, t *testing.T, rootTable string) {
	numGoroutines := 100
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		go func() {
			obj := mock.RandomPerson()

			ctx, cancel := getRaceContext()

			// will fail if rowsAff == 0
			testSave(ctx, o, t, obj)

			// will fail if rowsAff != 1
			testDelete(ctx, o, t, obj)
			cancel()

			wg.Done()
		}()

		wg.Add(1)
	}
	wg.Wait()
}
