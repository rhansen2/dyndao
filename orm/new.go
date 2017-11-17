package orm

import (
	"database/sql"

	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

// ORM is the primary object we expect the caller to operate on.
// Construct one with orm.New( ... ) and be on your merry way.
type ORM struct {
	sqlGen  *sg.SQLGenerator
	s       *schema.Schema
	RawConn *sql.DB

	// string is the table name that corresponds to a table in the schema. HookFunction
	BeforeCreateHooks map[string]HookFunction
	AfterCreateHooks  map[string]HookFunction

	BeforeUpdateHooks map[string]HookFunction
	AfterUpdateHooks  map[string]HookFunction

	BeforeDeleteHooks map[string]HookFunction
	AfterDeleteHooks  map[string]HookFunction
}

// GetSchema returns the ORM's active schema
func (o ORM) GetSchema() *schema.Schema {
	return o.s
}

func (o ORM) UseTracing() bool {
	return o.sqlGen.Tracing
}

// GetSQLGenerator returns the active SQLGenerator adapter
func (o ORM) GetSQLGenerator() *sg.SQLGenerator {
	return o.sqlGen
}

// New is the ORM constructor. It expects a SQL generator, JSON/SQL Schema object, and database connection.
func New(gen *sg.SQLGenerator, s *schema.Schema, db *sql.DB) ORM {
	o := ORM{sqlGen: gen, s: s, RawConn: db}

	o.BeforeCreateHooks = makeEmptyHookMap()
	o.AfterCreateHooks = makeEmptyHookMap()

	o.BeforeUpdateHooks = makeEmptyHookMap()
	o.AfterUpdateHooks = makeEmptyHookMap()

	o.BeforeDeleteHooks = makeEmptyHookMap()
	o.AfterDeleteHooks = makeEmptyHookMap()

	return o
}

