package mysql

import (
	sg "github.com/rbastic/dyndao/sqlgen"
)

// New shows off a sort of inheritance/composition-using-vtables approach.
// It receives the SQLGenerator composed by Core and then overrides any
// methods that it needs to. In some instances, this could be all methods,
// or hardly any.
func New(g *sg.SQLGenerator) *sg.SQLGenerator {
	g.IsMYSQL = true
	g.IsStringType = sg.FnIsStringType(IsStringType)
	g.IsNumberType = sg.FnIsNumberType(IsNumberType)
	g.IsFloatingType = sg.FnIsFloatingType(IsFloatingType)
	g.IsTimestampType = sg.FnIsTimestampType(IsTimestampType)
	g.IsLOBType = sg.FnIsLOBType(IsLOBType)
	g.RenderCreateColumn = sg.FnRenderCreateColumn(RenderCreateColumn)
	return g
}
