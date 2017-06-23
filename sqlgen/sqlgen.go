package sqlgen

import "errors"

type Generator struct {
	database string
	// TODO: functor setup
}

func New(db string) (Generator, error) {
	switch db {
	case "sqlite":
		//	fallthrough
		//case "oracle":
		return &sqlitegen.Generator{database: db}, nil
	default:
		return nil, errors.New("sqlgen: Unrecognized database type " + db)
	}
}
