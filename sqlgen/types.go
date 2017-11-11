package sqlgen

import (
	"database/sql"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
)

type FnBindingInsert func(g *SQLGenerator, sch *schema.Schema, table string, data map[string]interface{}) (string, []interface{}, error)
type FnBindingUpdate func(g *SQLGenerator, sch *schema.Schema, obj *object.Object) (string, []interface{}, []interface{}, error)
type FnBindingRetrieve func(g *SQLGenerator, sch *schema.Schema, obj *object.Object) (string, []string, []interface{}, error)
type FnBindingDelete func(g *SQLGenerator, sch *schema.Schema, obj *object.Object) (string, []interface{}, error)
type FnCreateTable func(g *SQLGenerator, sch *schema.Schema, table string) (string, error)
type FnDropTable func(name string) string
type FnRenderBindingValueWithInt func(f *schema.Column, i int) string
type FnRenderInsertValue func(bindI * int, f *schema.Column, value interface{}) (interface{}, error)
type FnIsStringType func(string) bool
type FnIsNumberType func(string) bool
type FnIsFloatingType func(string) bool
type FnIsTimestampType func(string) bool
type FnIsLOBType func(string) bool
type FnDynamicObjectSetter func(g *SQLGenerator, columnNames []string, columnPointers []interface{}, columnTypes []*sql.ColumnType, obj *object.Object) error
type FnMakeColumnPointers func(g *SQLGenerator, sliceLen int, columnTypes []*sql.ColumnType) ([]interface{}, error)

type FnRenderWhereClause func(g *SQLGenerator, schTable *schema.Table, obj *object.Object) (string, []interface{}, error)
type FnRenderUpdateWhereClause func(g *SQLGenerator, schTable *schema.Table, fieldsMap map[string]*schema.Column, obj *object.Object) (string, []interface{}, *int, error)

type FnCoreBindingInsert func(g *SQLGenerator, schTable *schema.Table, data map[string]interface{}, identityCol string, fieldsMap map[string]*schema.Column) ([]string, []string, []interface{})

type FnRenderCreateColumn func(g *SQLGenerator, f *schema.Column) string
type FnBindingInsertSQL func(schTable *schema.Table, tableName string, colNames []string, bindNames []string, identityCol string) string

// SQLGenerator is the 'vtable struct' that an ORM expects a SQL string
// generator to support.  While this does add an extra layer of indirection at
// runtime, it allows us to share common SQL idioms between implementations
// much more easily.
type SQLGenerator struct {
	Tracing                   bool
	FixLastInsertIDbug        bool
	// Necessary for ORM-level compatibility hacks
	IsSQLITE bool
	IsMYSQL bool
	IsORACLE bool
	IsMSSQL bool // MS SQL Server
	IsPOSTGRES bool // Postgre/Postgres? 's' or not?
	IsDB2 bool

	BindingInsert             FnBindingInsert
	BindingUpdate             FnBindingUpdate
	BindingRetrieve           FnBindingRetrieve
	BindingDelete             FnBindingDelete
	CreateTable               FnCreateTable
	RenderCreateColumn        FnRenderCreateColumn
	DropTable                 FnDropTable
	RenderBindingValueWithInt FnRenderBindingValueWithInt
	RenderInsertValue         FnRenderInsertValue

	IsStringType FnIsStringType

	IsNumberType    FnIsNumberType
	IsFloatingType  FnIsFloatingType
	IsTimestampType FnIsTimestampType
	IsLOBType       FnIsLOBType

	DynamicObjectSetter FnDynamicObjectSetter
	MakeColumnPointers  FnMakeColumnPointers

	RenderWhereClause       FnRenderWhereClause
	RenderUpdateWhereClause FnRenderUpdateWhereClause
	CoreBindingInsert       FnCoreBindingInsert
	BindingInsertSQL        FnBindingInsertSQL
}
