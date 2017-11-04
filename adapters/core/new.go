package core

import (
	sg "github.com/rbastic/dyndao/sqlgen"
	"os"
)

func New() *sg.SQLGenerator {
	g := new(sg.SQLGenerator)

	if os.Getenv("DB_TRACE") != "" {
		g.Tracing = true
	}

	g.CreateTable = sg.FnCreateTable(CreateTable)
	g.DropTable = sg.FnDropTable(DropTable)
	g.CoreBindingInsert = sg.FnCoreBindingInsert(CoreBindingInsert)
	g.BindingInsert = sg.FnBindingInsert(BindingInsert)
	g.BindingInsertSQL = sg.FnBindingInsertSQL(BindingInsertSQL)
	g.BindingRetrieve = sg.FnBindingRetrieve(BindingRetrieve)
	g.BindingUpdate = sg.FnBindingUpdate(BindingUpdate)
	g.BindingDelete = sg.FnBindingDelete(BindingDelete)
	g.RenderBindingValue = sg.FnRenderBindingValue(RenderBindingValue)
	g.RenderBindingValueWithInt = sg.FnRenderBindingValueWithInt(RenderBindingValueWithInt)
	g.RenderWhereClause = sg.FnRenderWhereClause(RenderWhereClause)
	g.RenderInsertValue = sg.FnRenderInsertValue(RenderInsertValue)
	g.RenderUpdateWhereClause = sg.FnRenderUpdateWhereClause(RenderUpdateWhereClause)
	g.DynamicObjectSetter = sg.FnDynamicObjectSetter(DynamicObjectSetter)
	g.MakeColumnPointers = sg.FnMakeColumnPointers(MakeColumnPointers)
	return g
}
