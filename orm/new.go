package orm

import (
	"database/sql"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

// Generator is the interface that an ORM expects a SQL string generator to support.
type Generator interface {
	BindingInsert(sch *schema.Schema, table string, data map[string]interface{}) (string, []interface{}, error)
	BindingUpdate(sch *schema.Schema, obj *object.Object) (string, []interface{}, []interface{}, error)
	BindingRetrieve(sch *schema.Schema, obj *object.Object) (string, []string, []interface{}, error)
	BindingDelete(sch *schema.Schema, obj *object.Object) (string, []interface{}, []interface{}, error)
	CreateTable(sch *schema.Schema, table string) (string, error)
	DropTable(name string) string

	IsStringType(string) bool
	IsNumberType(string) bool

	FixLastInsertIDbug() bool
}

// ORM is the abstraction that results from combining a sql generator, schema, and a database connection.
type ORM struct {
	sqlGen  Generator
	s       *schema.Schema
	RawConn *sql.DB
}

func (o ORM) GetSchema() *schema.Schema {
	return o.s
}

// New is the ORM constructor. It expects a SQL generator, JSON/SQL Schema object, and database connection.
func New(gen Generator, s *schema.Schema, db *sql.DB) ORM {
	o := ORM{sqlGen: gen, s: s, RawConn: db}
	return o
}
