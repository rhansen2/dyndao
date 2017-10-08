package orm

import (
	"database/sql"
)

// GetDB accepts a DSN and Driver as strings, returning a sql.DB and an error.
func GetDB(dsn, driver string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
