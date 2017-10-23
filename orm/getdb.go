package orm

import (
	"database/sql"
)

// GetDB accepts a Driver and DSN as strings, returning a sql.DB and an error.
func GetDB(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
