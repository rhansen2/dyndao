package oraclegen

import (
	sg "github.com/rbastic/dyndao/sqlgen"
)

// New shows off a sort of inheritance/composition-using-vtables approach.
// It receives the SQLGenerator composed by Core and then overrides any
// methods that it needs to. In some instances, this could be all methods,
// or hardly any.
func New(g * sg.SQLGenerator) *sg.SQLGenerator {
	// Oracle SQLGenerator uses Core for anything commented out.
	//g.BindingInsert = sg.FnBindingInsert(BindingInsert)
	//g.CoreBindingInsert = sg.FnCoreBindingInsert(CoreBindingInsert)
	g.BindingUpdate = sg.FnBindingUpdate(BindingUpdate)
	g.BindingRetrieve = sg.FnBindingRetrieve(BindingRetrieve)
	g.BindingDelete = sg.FnBindingDelete(BindingDelete)

	g.RenderBindingValue = sg.FnRenderBindingValue(RenderBindingValue)
	//g.CreateTable = sg.FnCreateTable(CreateTable)
	g.DropTable = sg.FnDropTable(DropTable)
	g.FixLastInsertIDbug = sg.FnFixLastInsertIDbug(FixLastInsertIDbug)
	g.RenderBindingValueWithInt = sg.FnRenderBindingValueWithInt(RenderBindingValueWithInt)
	g.IsStringType = sg.FnIsStringType(IsStringType)
	g.IsNumberType = sg.FnIsNumberType(IsNumberType)
	g.IsFloatingType = sg.FnIsFloatingType(IsFloatingType)
	g.IsTimestampType = sg.FnIsTimestampType(IsTimestampType)
	g.IsLOBType = sg.FnIsLOBType(IsLOBType)
	g.DynamicObjectSetter = sg.FnDynamicObjectSetter(DynamicObjectSetter)
	g.MakeColumnPointers = sg.FnMakeColumnPointers(MakeColumnPointers)
	g.RenderWhereClause = sg.FnRenderWhereClause(RenderWhereClause)
	g.RenderUpdateWhereClause = sg.FnRenderUpdateWhereClause(RenderUpdateWhereClause)
	g.RenderCreateField = sg.FnRenderCreateField(RenderCreateField)
	g.RenderInsertValue = sg.FnRenderInsertValue(RenderInsertValue)
	return g
}
