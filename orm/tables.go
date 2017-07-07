package orm

import (
	"context"
	"database/sql"
	"os"
	// TODO: Use log15 instead of fmt?
	"fmt"

	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/schema"
)

func (o ORM) CreateTables() error {
	for tName := range o.s.Tables {
		err := o.CreateTable(o.s, tName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o ORM) DropTables() error {
	for tName := range o.s.Tables {
		err := o.DropTable(tName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o ORM) CreateTable(sch *schema.Schema, tableName string) error {
	sqlStr, err := o.sqlGen.CreateTable(sch, tableName)
	if err != nil {
		return err
	}

	debug := os.Getenv("DEBUG")
	if debug != "" {
		// Currently, DEBUG is either on or off.
		fmt.Println("CreateTable:", sqlStr)
	}

	_, err = prepareAndExecSQL(o.RawConn, sqlStr)
	if err != nil {
		return errors.Wrap(err, "CreateTable")
	}
	return nil
}

func (o ORM) DropTable(tableName string) error {
	sqlStr := o.sqlGen.DropTable(tableName)
	_, err := prepareAndExecSQL(o.RawConn, sqlStr)
	if err != nil {
		return errors.Wrap(err, "CreateTable")
	}
	return nil
}

func prepareAndExecSQL(db *sql.DB, sqlStr string) (sql.Result, error) {
	stmt, err := db.PrepareContext(context.TODO(), sqlStr)
	defer stmt.Close()
	r, err := stmt.ExecContext(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "prepareAndExecSQL (" + sqlStr +")")
	}
	return r, nil
}