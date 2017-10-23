// Package infoschema is a schema parser for the information_schema
// standard.
//
// By supporting the information schema, it is possible to directly
// support loading schemas dynamically for at least the following
// databases:
//
// Microsoft SQL Server, MySQL, PostgreSQL, InterSystems Caché
// H2 Database, HSQLDB, MariaDB
//

/*
	TODO: foreign key identification
	TODO: Indexes
	TODO: Constraints
	TODO: interface type for infoschema package (so that we
	have an interface to implement an oracle-alike version for,
	since Oracle doesn't support info-schema... and requires some
	open source package to implement it otherwise, which, you know,
	who has time to install things like that!)

	TODO: oracle equivalent that implements aforementioned interface
*/

package infoschema

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rbastic/dyndao/schema"
)

var (
	INFO_TABLES = "information_schema.tables"
)

func getTableNamesSQL(db string) string {
	return fmt.Sprintf(`
SELECT DISTINCT TABLE_NAME 
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA='%s'
	`, db)
}

// LoadSchema loads the entire schema and configures the essential
// fields to be all columns in the table.
func LoadSchema(ctx context.Context, db *sql.DB, dbName string) (*schema.Schema, error) {
	sch, err := ParseSchema(ctx, db, dbName)
	if err != nil {
		return nil, err
	}
	err = ParseTables(ctx, db, dbName, sch)
	if err != nil {
		return nil, err
	}
	SetDefaultEssentialFields(sch)
	return sch, nil
}

// ParseSchema does a preliminary load of the schema, reading in all
// table names and populating default schema.Table structures.
func ParseSchema(ctx context.Context, db *sql.DB, dbName string) (*schema.Schema, error) {
	rows, err := db.QueryContext(ctx, getTableNamesSQL(dbName))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()

	sch := schema.DefaultSchema()
	for rows.Next() {
		var tblName string
		err := rows.Scan(&tblName)
		if err != nil {
			return nil, err
		}
		sch.Tables[tblName] = schema.DefaultTable()
	}
	err = rows.Err()
	return sch, err
}

func getColumnMetaSQL(db string) string {
	return fmt.Sprintf(`
SELECT TABLE_NAME, COLUMN_NAME, DATA_TYPE, COLUMN_DEFAULT, IS_NULLABLE, COLUMN_KEY, EXTRA
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA='%s'
ORDER BY TABLE_NAME
	`, db)
}

// ParseTables loads all potential column information from a given schema into the relevant tables.
func ParseTables(ctx context.Context, db *sql.DB, dbName string, sch *schema.Schema) error {
	rows, err := db.QueryContext(ctx, getColumnMetaSQL(dbName))
	if err != nil {
		return err
	}
	defer func() {
		err = rows.Close()
	}()

	for rows.Next() {
		var tblName string
		var colName sql.NullString
		var dataType string
		var columnDefault sql.NullString
		var isNullable string
		var columnKey string
		var extra string

		err := rows.Scan(&tblName, &colName, &dataType, &columnDefault, &isNullable, &columnKey, &extra)
		if err != nil {
			return err
		}

		// Mutates the schema.Table for the given tblName and colName
		setTableCol(sch, tblName, colName, dataType, columnDefault, isNullable, columnKey, extra)
	}

	err = rows.Err()
	return err
}

func setTableCol(sch *schema.Schema, tblName string, colName sql.NullString, dataType string, colDefault sql.NullString, isNullable string, columnKey string, extra string) {
	tbl := sch.Tables[tblName]
	tbl.Name = tblName

	df := schema.DefaultField()
	df.Name = colName.String
	df.DBType = dataType
	df.DefaultValue = colDefault.String

	isNullBool := false
	if isNullable == "YES" {
		isNullBool = true
	}
	df.AllowNull = isNullBool

	// TODO: is this a standard or just MySQL that supports
	// this?
	isIdentity := false
	if extra == "auto_increment" && columnKey == "PRI" {
		isIdentity = true
	}
	df.IsIdentity = isIdentity

	// TODO: IsNumber, need a SQL generator for that, unless we deprecate
	// IsNumber, I think.  See issue #49 on github.
	tbl.Fields[colName.String] = df
}

// SetDefaultEssentialFields configures the EssentialFields
// for each schema.Table to be the entire list of field names.
func SetDefaultEssentialFields(sch *schema.Schema) {
	for _, tbl := range sch.Tables {
		numf := len(tbl.Fields)
		tbl.EssentialFields = make([]string, numf)

		i := 0
		for _, v := range tbl.Fields {
			tbl.EssentialFields[i] = v.Name
			i++
		}
	}
}
