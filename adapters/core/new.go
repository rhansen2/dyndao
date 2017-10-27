package core

import (
	sg "github.com/rbastic/dyndao/sqlgen"
)

func New() *sg.SQLGenerator {
	g := new(sg.SQLGenerator)
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
	return g
}