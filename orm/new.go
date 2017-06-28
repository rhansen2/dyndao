package orm

import (
	"database/sql"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

type Generator interface {
	BindingInsert(sch *schema.Schema, table string, data map[string]interface{}) (string, []interface{}, error)
	BindingUpdate(sch *schema.Schema, obj *object.Object) (string, []interface{}, []interface{}, error)
	BindingRetrieve(sch *schema.Schema, obj *object.Object) (string, []interface{}, error)
	CreateTable(sch *schema.Schema, table string) (string, error)
	DropTable(name string) string
}

type ORM struct {
	sqlGen  Generator
	s       *schema.Schema
	RawConn *sql.DB
}

func New(gen Generator, s *schema.Schema, db *sql.DB) ORM {
	o := ORM{sqlGen: gen, s: s, RawConn: db}
	return o
}
