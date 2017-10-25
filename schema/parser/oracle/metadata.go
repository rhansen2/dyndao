// Package oracle is a schema parser for the Oracle metadata.
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

package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/rbastic/dyndao/schema"
)

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
	SetDefaultEssentialColumns(sch)
	return sch, nil
}

// db string here is what Oracle calls the 'TABLESPACE'.
func getTableNamesSQL(db string) string {
	return fmt.Sprintf(`SELECT DISTINCT(TABLE_NAME) FROM ALL_TABLES WHERE TABLESPACE_NAME=:TBN0 AND STATUS = :STAT1`)
}

// ParseSchema does a preliminary load of the schema, reading in all
// table names and populating default schema.Table structures.
func ParseSchema(ctx context.Context, db *sql.DB, dbName string) (*schema.Schema, error) {
	sqlStr := getTableNamesSQL(dbName)
	if os.Getenv("DB_TRACE") != "" {
		fmt.Println("dyndao: ParseSchema SQL: ", sqlStr)
	}

	rows, err := db.QueryContext(ctx, sqlStr, dbName, "VALID")
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

func getColumnMetaSQL(db string, tblName string) string {
	return fmt.Sprintf(`DESCRIBE '%s'.%s`, db, tblName)
}

// ParseTables loads all potential column information from a given schema into the relevant tables.
func ParseTables(ctx context.Context, db *sql.DB, dbName string, sch *schema.Schema) error {
	for _, tbl := range sch.Tables {
		rows, err := db.QueryContext(ctx, getColumnMetaSQL(dbName, tbl.Name))
		if err != nil {
			return err
		}

		defer func() {
			err = rows.Close()
		}()

		for rows.Next() {
			var colName sql.NullString
			var isNullable string
			var dataType string

			err := rows.Scan(&colName, &isNullable, &dataType)
			if err != nil {
				return err
			}

			// Mutates the schema.Table for the given tblName and colName
			setTableCol(sch, tbl.Name, colName, dataType, isNullable)
		}

		err = rows.Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func setTableCol(sch *schema.Schema, tblName string, colName sql.NullString, dataType string, isNullable string) {
	tbl := sch.Tables[tblName]
	tbl.Name = tblName

	df := schema.DefaultColumn()
	df.Name = colName.String
	df.DBType = dataType

	isNullBool := false
	if isNullable == "NULL" {
		isNullBool = true
	}
	df.AllowNull = isNullBool

	tbl.Columns[colName.String] = df
}

// SetDefaultEssentialColumns configures the EssentialColumns
// for each schema.Table to be the entire list of field names.
func SetDefaultEssentialColumns(sch *schema.Schema) {
	for _, tbl := range sch.Tables {
		numf := len(tbl.Columns)
		tbl.EssentialColumns = make([]string, numf)

		i := 0
		for _, v := range tbl.Columns {
			tbl.EssentialColumns[i] = v.Name
			i++
		}
	}
}
