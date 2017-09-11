package orm

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

func prepareAndExecSQL(db *sql.DB, sqlStr string) (sql.Result, error) {
	stmt, err := db.PrepareContext(context.TODO(), sqlStr)
	if err != nil {
		return nil, errors.Wrap(err, "prepareAndExecSQL/PrepareContext ("+sqlStr+")")
	}
	defer stmt.Close()
	r, err := stmt.ExecContext(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "prepareAndExecSQL/ExecContext ("+sqlStr+")")
	}
	return r, nil
}
