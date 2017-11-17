package orm

import (
	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/dyndao/object"
)

// HookFunction is the function type for declaring a software-based trigger, which we
// refer to as a 'hook function'.
type HookFunction func(*schema.Schema, *object.Object) error

// See new.go for the ORM type, which contains the
// declarations for BeforeCreateHooks, AfterCreateHooks, etc.

func makeEmptyHookMap() map[string]HookFunction {
	return make(map[string]HookFunction)
}

// Software trigger functions
func (o *ORM) CallBeforeCreateHookIfNeeded(obj *object.Object) error {
	hookFunc, ok := o.BeforeCreateHooks[o.s.GetTableName(obj.Type)]
	if ok {
		return hookFunc(o.s, obj)
	}
	return nil
}

func (o *ORM) CallAfterCreateHookIfNeeded(obj *object.Object) error {
	hookFunc, ok := o.AfterCreateHooks[o.s.GetTableName(obj.Type)]
	if ok {
		return hookFunc(o.s, obj)
	}
	return nil

}

func (o *ORM) CallBeforeUpdateHookIfNeeded(obj *object.Object) error {
	hookFunc, ok := o.BeforeUpdateHooks[o.s.GetTableName(obj.Type)]
	if ok {
		return hookFunc(o.s, obj)
	}
	return nil
}

func (o *ORM) CallAfterUpdateHookIfNeeded(obj *object.Object) error {
	hookFunc, ok := o.AfterUpdateHooks[o.s.GetTableName(obj.Type)]
	if ok {
		return hookFunc(o.s, obj)
	}
	return nil
}

func (o *ORM) CallBeforeDeleteHookIfNeeded(obj *object.Object) error {
	hookFunc, ok := o.BeforeDeleteHooks[o.s.GetTableName(obj.Type)]
	if ok {
		return hookFunc(o.s, obj)
	}
	return nil
}

func (o *ORM) CallAfterDeleteHookIfNeeded(obj *object.Object) error {
	hookFunc, ok := o.AfterDeleteHooks[o.s.GetTableName(obj.Type)]
	if ok {
		return hookFunc(o.s, obj)
	}
	return nil
}
