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
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rbastic/dyndao/schema"
)

// LoadSchema loads the entire schema and configures the essential
// fields to be all columns in the table.
func LoadSchema(ctx context.Context, db *sql.DB, dbName string) (*schema.Schema, error) {
	sch, err := ParseSchema(ctx, db, dbName)
	if err != nil {
		return nil, errors.Wrap(err, "LoadSchema")
	}
	err = ParseTables(ctx, db, dbName, sch)
	if err != nil {
		return nil, errors.Wrap(err, "ParseTables")
	}
	SetDefaultEssentialColumns(sch)
	return sch, nil
}

// db string here is what Oracle calls the 'TABLESPACE'.
func getTableNamesSQL(db string) string {
	return fmt.Sprintf("select table_name from dba_tables where owner like '%s'", strings.ToUpper(db))
}

// ParseSchema does a preliminary load of the schema, reading in all
// table names and populating default schema.Table structures.
func ParseSchema(ctx context.Context, db *sql.DB, dbName string) (*schema.Schema, error) {
	sqlStr := getTableNamesSQL(dbName)
	if os.Getenv("DB_TRACE") != "" {
		fmt.Printf("dyndao: ParseSchema SQL: [%s] bindArgs: [%v] [%v]\n", sqlStr, dbName, "VALID")
	}

	//rows, err := db.QueryContext(ctx, sqlStr, dbName, "VALID")
	rows, err := db.QueryContext(ctx, sqlStr)
	if err != nil {
		return nil, errors.Wrap(err, "ParseSchema/QueryContext")
	}
	defer func() {
		err = rows.Close()
	}()

	sch := schema.DefaultSchema()
	for rows.Next() {
		var tblName string

		err := rows.Scan(&tblName)
		if err != nil {
			wrapMsg := "ParseSchema/rows.Scan(\"" + tblName + "\")"
			return nil, errors.Wrap(err, wrapMsg)
		}

		if shouldSkipParsingTable(tblName) {
			continue
		}

		schTbl := schema.DefaultTable()
		schTbl.Name = tblName
		sch.Tables[tblName] = schTbl
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "rows.Err()")
	}

	return sch, err
}

func shouldSkipParsingTable(tblName string) bool {
	// -HACK- skip any tables with $ in the name.
	if strings.Contains(tblName, "$") {
		return true
	}

	return false
}

func getColumnMetaSQL(db string, tblName string) string {
	if db == "" {
		panic("getColumnMetaSQL expected non-empty db parameter")
	}
	if tblName == "" {
		panic("getColumnMetaSQL expected non-empty tblName parameter")

	}
	// RPAD(COLUMN_NAME,30)||' '||DATA_TYPE||'('||DATA_LENGTH||')' as descr

	// TODO: desc all_tab_cols from within system cli -- implement more
	return fmt.Sprintf(`
 select COLUMN_NAME, DATA_TYPE, DATA_LENGTH, NULLABLE, IDENTITY_COLUMN
 FROM all_tab_cols
 WHERE TABLE_NAME ='%s'  and owner= '%s'
`, tblName, db)
}

// ParseTables loads all potential column information from a given schema into the relevant tables.
func ParseTables(ctx context.Context, db *sql.DB, dbName string, sch *schema.Schema) error {
	for _, tbl := range sch.Tables {
		metasql := getColumnMetaSQL(dbName, tbl.Name)
		if os.Getenv("DB_TRACE") != "" {
			fmt.Printf("dyndao: ParseTables getColumnMetaSQL: %s\n", metasql)
		}
		rows, err := db.QueryContext(ctx, metasql)
		if err != nil {
			return errors.Wrap(err, "QueryContext")
		}

		defer func() {
			err = rows.Close()
		}()

		for rows.Next() {
			var colName sql.NullString
			var dataType string
			var dataLength string
			var isNullable string
			var identityCol string

			err := rows.Scan(&colName, &dataType, &dataLength, &isNullable, &identityCol)
			if err != nil {
				return errors.Wrap(err, "rows.Scan()")
			}

			// Mutates the schema.Table for the given tblName and colName
			setTableCol(sch, tbl.Name, colName, dataType, dataLength, isNullable, identityCol)
		}

		err = rows.Err()
		if err != nil {
			return errors.Wrap(err, "rows.Err()")
		}
	}
	return nil
}

func setTableCol(sch *schema.Schema, tblName string, colName sql.NullString, dataType string, dataLength string, isNullable string, identityCol string) {
	tbl := sch.Tables[tblName]
	tbl.Name = tblName

	df := schema.DefaultColumn()
	df.Name = colName.String
	df.DBType = dataType

	isNullBool := false
	if isNullable == "Y" {
		isNullBool = true
	}

	//	In case you get curious.
	//	fmt.Println("setTableCol: dataType:[", dataType, "]")
	//	fmt.Println("setTableCol: isNullable:[", isNullable, "] identityCol:[", identityCol, "]")
	//	fmt.Println("setTableCol: dataLength:[", dataLength, "]")

	df.AllowNull = isNullBool

	dli, err := strconv.Atoi(dataLength)
	if err != nil {
		panic(err)
	}
	df.Length = dli

	df.IsIdentity = true
	if identityCol == "NO" {
		df.IsIdentity = false
	}

	// TODO: Consider adding a 'SetDataType' function to dyndao that lets
	// us better work with situations like this
	if dataType == "NUMBER" {
		df.IsNumber = true
	}

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
