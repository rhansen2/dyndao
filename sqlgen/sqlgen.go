package sqlgen

import "errors"
import "github.com/rbastic/dyndao/sqlgen/sqlitegen"
import "github.com/rbastic/dyndao/schema"

func New(db string, name string, sch * schema.Schema) (interface{}, error) {
	switch db {
	case "sqlite":
		//	fallthrough
		//case "oracle":
		// TODO: fix testName as a parameter
		return sqlitegen.New( db, "testName", sch ), nil
	default:
		return nil, errors.New("sqlgen: Unrecognized database type " + db)
	}
}
