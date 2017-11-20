package orm

import (
	"os"
	// TODO: Use log15 instead of fmt?
	"fmt"

	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/schema"
)

// CreateTables executes a CreateTable operation for every table specified in
// the schema.
func (o *ORM) CreateTables(ctx context.Context) error {
	for tName := range o.s.Tables {
		err := o.CreateTable(ctx, o.s, tName)
		if err != nil {
			return err
		}
	}

	return nil
}

// DropTables executes a DropTable operation for every table specified in the
// schema.
func (o *ORM) DropTables(ctx context.Context) error {
	for tName := range o.s.Tables {
		err := o.DropTable(ctx, tName)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateTable will execute a CreateTable operation for the specified table in
// a given schema.
func (o *ORM) CreateTable(ctx context.Context, sch *schema.Schema, tableName string) error {
	sqlStr, err := o.sqlGen.CreateTable(o.sqlGen, sch, tableName)
	if err != nil {
		return err
	}

	debug := os.Getenv("DB_TRACE")
	if debug != "" {
		// Currently, DEBUG is either on or off.
		fmt.Println("CreateTable:", sqlStr)
	}

	_, err = prepareAndExecSQL(ctx, o.RawConn, sqlStr)
	if err != nil {
		return errors.Wrap(err, "CreateTable")
	}
	return nil
}

// DropTable will execute a DropTable operation for the specified table in
// a given schema.
func (o *ORM) DropTable(ctx context.Context, tableName string) error {
	sqlStr := o.sqlGen.DropTable(tableName)
	_, err := prepareAndExecSQL(ctx, o.RawConn, sqlStr)
	if err != nil {
		return errors.Wrap(err, "DropTable")
	}
	return nil
}

func prepareAndExecSQL(ctx context.Context, db *sql.DB, sqlStr string) (sql.Result, error) {
	stmt, err := db.PrepareContext(ctx, sqlStr)
	if err != nil {
		return nil, errors.Wrap(err, "prepareAndExecSQL/PrepareContext ("+sqlStr+")")
	}
	defer func() {
		stmtErr := stmt.Close()
		if stmtErr != nil {
			fmt.Println(stmtErr) // TODO: logging implementation
		}
	}()
	r, err := stmt.ExecContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "prepareAndExecSQL/ExecContext ("+sqlStr+")")
	}
	return r, nil
}
