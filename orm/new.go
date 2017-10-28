package orm

import (
	"database/sql"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

// HookFunction is the function type for declaring a software-based trigger, which we
// refer to as a 'hook function'.
type HookFunction func(*schema.Schema, *object.Object) error

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

func makeEmptyHookMap() map[string]HookFunction {
	return make(map[string]HookFunction)
}

// New is the ORM constructor. It expects a SQL generator, JSON/SQL Schema object, and database connection.
func New(gen *sg.SQLGenerator, s *schema.Schema, db *sql.DB) ORM {
	o := ORM{sqlGen: gen, s: s, RawConn: db}

	o.BeforeCreateHooks = makeEmptyHookMap()
	o.AfterCreateHooks = makeEmptyHookMap()

	o.BeforeUpdateHooks = makeEmptyHookMap()
	o.AfterUpdateHooks = makeEmptyHookMap()

	return o
}

// Software trigger functions

// CallBeforeCreateHookIfNeeded will call the necessary BeforeCreate triggers for a given
// object if they are set.
func (o *ORM) CallBeforeCreateHookIfNeeded(obj *object.Object) error {
	hookFunc, ok := o.BeforeCreateHooks[o.s.GetTableName(obj.Type)]
	if ok {
		return hookFunc(o.s, obj)
	}
	return nil
}

// CallAfterCreateHookIfNeeded will call the necessary AfterCreate triggers for a given
// object if they are set.
func (o *ORM) CallAfterCreateHookIfNeeded(obj *object.Object) error {
	hookFunc, ok := o.AfterCreateHooks[o.s.GetTableName(obj.Type)]
	if ok {
		return hookFunc(o.s, obj)
	}
	return nil

}

// CallBeforeUpdateHookIfNeeded will call the necessary BeforeUpdate triggers for a given
// object if they are set.
func (o *ORM) CallBeforeUpdateHookIfNeeded(obj *object.Object) error {
	hookFunc, ok := o.BeforeUpdateHooks[o.s.GetTableName(obj.Type)]
	if ok {
		return hookFunc(o.s, obj)
	}
	return nil
}

// CallAfterUpdateHookIfNeeded will call the necessary AfterUpdate triggers for
// a given object if they are set.
func (o *ORM) CallAfterUpdateHookIfNeeded(obj *object.Object) error {
	hookFunc, ok := o.AfterUpdateHooks[o.s.GetTableName(obj.Type)]
	if ok {
		return hookFunc(o.s, obj)
	}
	return nil
}
