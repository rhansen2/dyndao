// Package tests is a set of tests that put the various components together and
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
	sch := schema.MockBasicSchema()
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

func TestSaveBasicObject(t *testing.T) {
	sch := schema.MockBasicSchema()
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

	{
		fmt.Println("Saving Ryan")
		rowsAff, err := o.SaveObject(context.TODO(), nil, obj)
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
		t.Fatalf("PersonID has the wrong value, has value %d", personID)
	}
	fmt.Println("PersonID=", personID)

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
		if latestJoe.Get("PersonID") != 1 || latestJoe.Get("Name") != "Joe" {
			t.Fatal("latestJoe does not match expectations")
		}
	}

}

func TestDropTables(t *testing.T) {
	sch := schema.MockBasicSchema()
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
