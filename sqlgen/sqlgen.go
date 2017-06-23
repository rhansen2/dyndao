package sqlgen

import "errors"
import "github.com/rbastic/dyndao/sqlgen/sqlitegen"

func New(db string) (interface{}, error) {
	switch db {
	case "sqlite":
		//	fallthrough
		//case "oracle":
		return &sqlitegen.Generator{Database: db}, nil
	default:
		return nil, errors.New("sqlgen: Unrecognized database type " + db)
	}
}
