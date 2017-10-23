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
	BindingDelete(sch *schema.Schema, obj *object.Object) (string, []interface{}, error)
	CreateTable(sch *schema.Schema, table string) (string, error)
	DropTable(name string) string

	RenderBindingValue(f *schema.Field) string
	RenderBindingValueWithInt(f *schema.Field, i int64) string

	IsStringType(string) bool
	IsNumberType(string) bool
	IsFloatingType(string) bool
	IsTimestampType(string) bool
	IsLOBType(string) bool

	// This option will turn MODE_LAST_INSERT_ID off? Start naming these
	// things all mode? Same with FixLastInsertIDbug()?
	FixLastInsertIDbug() bool

	DynamicObjectSetter(columnNames []string, columnPointers []interface{}, columnTypes []*sql.ColumnType, obj *object.Object) error
	MakeColumnPointers(sliceLen int, columnTypes []*sql.ColumnType) ([]interface{}, error)
}

// HookFunction is the function type for declaring a software-based trigger, which we
// refer to as a 'hook function'.
type HookFunction func(*schema.Schema, *object.Object) error

// ORM is the primary object we expect the caller to operate on.
// Construct one with orm.New( ... ) and be on your merry way.
type ORM struct {
	sqlGen  Generator
	s       *schema.Schema
	RawConn *sql.DB

	// string is the table name that corresponds to a table in the schema. HookFunction
	BeforeCreateHooks map[string]HookFunction
	AfterCreateHooks  map[string]HookFunction

	BeforeUpdateHooks map[string]HookFunction
	AfterUpdateHooks  map[string]HookFunction
}

// GetSchema returns the current schema object that is stored within
// a given ORM object.
func (o ORM) GetSchema() *schema.Schema {
	return o.s
}

// GetGenerator returns the current sql generator object that is stored within a given ORM object.
func (o ORM) GetGenerator() Generator {
	return o.sqlGen
}

func makeEmptyHookMap() map[string]HookFunction {
	return make(map[string]HookFunction)
}

// New is the ORM constructor. It expects a SQL generator, JSON/SQL Schema object, and database connection.
func New(gen Generator, s *schema.Schema, db *sql.DB) ORM {
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
