package cockroach

import (
	//_ "github.com/denisenkom/go-postgresdb"
	postgre "github.com/rbastic/dyndao/adapters/postgres"
	sg "github.com/rbastic/dyndao/sqlgen"
)

// New shows off a sort of inheritance/composition-using-vtables approach.
// It receives the SQLGenerator composed by Core and then overrides any
// methods that it needs to. In some instances, this could be all methods,
// or hardly any.
func New(g *sg.SQLGenerator) *sg.SQLGenerator {
	g.IsPOSTGRES = true
	g.FixLastInsertIDbug = true
	g.IsStringType = sg.FnIsStringType(postgre.IsStringType)
	g.IsNumberType = sg.FnIsNumberType(postgre.IsNumberType)
	g.IsFloatingType = sg.FnIsFloatingType(postgre.IsFloatingType)
	g.IsTimestampType = sg.FnIsTimestampType(postgre.IsTimestampType)
	g.IsLOBType = sg.FnIsLOBType(postgre.IsLOBType)
	g.BindingInsertSQL = sg.FnBindingInsertSQL(postgre.BindingInsertSQL)
	g.RenderCreateColumn = sg.FnRenderCreateColumn(postgre.RenderCreateColumn)
	g.RenderBindingValueWithInt = sg.FnRenderBindingValueWithInt(postgre.RenderBindingValueWithInt)
	g.BindingUpdate = sg.FnBindingUpdate(postgre.BindingUpdate)
	g.DynamicObjectSetter = sg.FnDynamicObjectSetter(postgre.DynamicObjectSetter)
	g.MakeColumnPointers = sg.FnMakeColumnPointers(postgre.MakeColumnPointers)
	return g
}
