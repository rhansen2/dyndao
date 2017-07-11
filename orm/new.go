package orm

// The ORM package is designed to tie everything together: a database connection, schema,
// relevant objects, etc. The current design is a WIP. While not finished, it is serviceable
// and can be used effectively.

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
	BindingDelete(sch *schema.Schema, obj *object.Object) (string, []interface{}, error)
	CreateTable(sch *schema.Schema, table string) (string, error)
	DropTable(name string) string

	IsStringType(string) bool
	IsNumberType(string) bool
	// TODO: Name supplying our primary key as CallerSuppliesPrimaryKey?
	// This option will turn MODE_LAST_INSERT_ID off? Start naming these
	// things all mode? Same with FixLastInsertIDbug()?
	FixLastInsertIDbug() bool
	CallerSuppliesPrimaryKey() bool
}

// ORM is the primary object we expect the caller to operate on.
// Construct one with orm.New( ... ) and be on your merry way.
type ORM struct {
	sqlGen  Generator
	s       *schema.Schema
	RawConn *sql.DB
}

// GetSchema returns the current schema object that is stored within
// a given ORM object.
func (o ORM) GetSchema() *schema.Schema {
	return o.s
}

// New is the ORM constructor. It expects a SQL generator, JSON/SQL Schema object, and database connection.
func New(gen Generator, s *schema.Schema, db *sql.DB) ORM {
	o := ORM{sqlGen: gen, s: s, RawConn: db}
	return o
}
