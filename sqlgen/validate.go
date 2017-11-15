package sqlgen

func PanicIfInvalid(g *SQLGenerator) {
	if g.BindingInsert == nil {
		panic("dyndao: vtable BindingInsert is nil")
	}
	if g.BindingUpdate == nil {
		panic("dyndao: vtable BindingUpdate is nil")
	}
	if g.BindingRetrieve == nil {
		panic("dyndao: vtable BindingRetrieve is nil")
	}
	if g.BindingDelete == nil {
		panic("dyndao: vtable BindingDelete is nil")
	}
	if g.CreateTable == nil {
		panic("dyndao: vtable CreateTable is nil")
	}
	if g.DropTable == nil {
		panic("dyndao: vtable DropTable is nil")
	}
	if g.RenderBindingValueWithInt == nil {
		panic("dyndao: vtable RenderBindingValueWithInt is nil")
	}
	if g.IsStringType == nil {
		panic("dyndao: vtable IsStringType is nil")
	}
	if g.IsNumberType == nil {
		panic("dyndao: vtable IsNumberType is nil")
	}
	if g.IsFloatingType == nil {
		panic("dyndao: vtable IsFloatingType is nil")
	}
	if g.IsTimestampType == nil {
		panic("dyndao: vtable IsTimestampType is nil")
	}
	if g.IsLOBType == nil {
		panic("dyndao: vtable IsLOBType is nil")
	}
	if g.DynamicObjectSetter == nil {
		panic("dyndao: vtable DynamicObjectSetter is nil")
	}
	if g.MakeColumnPointers == nil {
		panic("dyndao: vtable MakeColumnPointers is nil")
	}
	if g.RenderWhereClause == nil {
		panic("dyndao: vtable RenderWhereClause is nil")
	}
	if g.RenderUpdateWhereClause == nil {
		panic("dyndao: vtable RenderUpdateWhereClause is nil")
	}
	if g.CoreBindingInsert == nil {
		panic("dyndao: vtable CoreBindingInsert is nil")
	}
	if g.RenderCreateColumn == nil {
		panic("dyndao: vtable RenderCreateColumn is nil")
	}
	if g.RenderInsertValue == nil {
		panic("dyndao: vtable RenderInsertValue is nil")
	}
	if g.BindingInsertSQL == nil {
		panic("dyndao: vtable BindingInsertSQL is nil")
	}
}
