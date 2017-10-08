// Package object is an abstract data record which tracks state changes.  It's
// meant to make it easier to map key-value records into ORM / RDBMS systems.
// The state change tracking can be useful when the values of primary keys need
// to be updated. (Changing a foreign key on a table with a composite key, for
// example)
package object

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	// ErrValueWasNil is returned by the Get* family of functions in
	// situations where a value must be returned, but we would like to
	// signal that it was originally nil.
	ErrValueWasNil = errors.New("object: value was nil")
)

// Array is our 'object array container' to assist with a couple of instances
// where slices are needed.
type Array []*Object

// Object struct encapsulates our key-value pairs (KV) and a single-item
// per-key history of the previous value stored for a given key
// (ChangedFields).  We also store any instances of 'child records' which may
// be relevant (for instance, when saving with nested transactions).  'saved'
// is used to track the internal state of whether an object was recently
// retrieved or remapped from internal database state.
type Object struct {
	Type          string
	KV            map[string]interface{}
	HiddenKV      map[string]interface{}
	ChangedFields map[string]interface{}
	Children      map[string]Array
	saved         bool
}

// New is an empty constructor
func New(objType string) *Object {
	return &Object{Type: objType, KV: makeEmptyMap(), HiddenKV: nil, ChangedFields: makeEmptyMap(), Children: makeEmptyChildrenMap(), saved: false}
}

// MakeArray will construct an array of length 'size'
func MakeArray(size int) Array {
	objAry := make(Array, size)
	return objAry
}

// NewArray is our single-object array constructor
func NewArray(val *Object) Array {
	objAry := make(Array, 1)
	objAry[0] = val
	return objAry
}

// NewEmptyArray is an empty array constructor
func NewEmptyArray() Array {
	objAry := make(Array, 0)
	return objAry
}

// just a bit of shorthand for the methods below
func makeEmptyMap() map[string]interface{} {
	return make(map[string]interface{})
}

func makeEmptyChildrenMap() map[string]Array {
	return make(map[string]Array)
}

func (o *Object) MakeHiddenKVIfNeeded() {
	if o.HiddenKV == nil {
		o.HiddenKV = makeEmptyMap()
	}
}

// Get is the most basic accessor, for cases
// that may not be handled by other methods
func (o Object) Get(k string) interface{} {
	return o.KV[k]
}

// GetWithFlag is the second most basic accessor, for cases
// that may not be handled by other methods
func (o Object) GetWithFlag(k string) (interface{}, bool) {
	v, ok := o.KV[k]
	return v, ok
}

// GetBool is a safe, typed bool accessor
func (o Object) GetBool(k string) (bool, bool) {
	v, ok := o.KV[k].(bool)
	return v, ok
}

// GetString is a safe, typed string accessor
func (o Object) GetString(k string) (string, bool) {
	v, ok := o.KV[k].(string)
	return v, ok
}

// GetInt is a safe, typed int64 accessor
func (o Object) GetInt(k string) (int64, bool) {
	v, ok := o.KV[k].(int64)
	return v, ok
}

// GetFloat is a safe, typed float64 accessor
func (o Object) GetFloat(k string) (float64, bool) {
	v, ok := o.KV[k].(float64)
	return v, ok
}

// HiddenGetStringAlways is a safe, typed string accessor for the Hidden KV. It
// will force conversion away from float64, int64, uint64, string, and nil
// values. Nils and unrecognized values are marked as an error (nil values will
// return 0 and ErrValueWasNil)
func (o Object) HiddenGetStringAlways(k string) (string, error) {
	switch v := o.HiddenKV[k].(type) {
	case float64:
		fl := o.HiddenKV[k].(float64)
		return fmt.Sprintf("%f", fl), nil
	case int64:
		fl := o.HiddenKV[k].(int64)
		return fmt.Sprintf("%d", fl), nil
	case uint64:
		fl := o.HiddenKV[k].(uint64)
		return fmt.Sprintf("%d", fl), nil
	case string:
		fl := o.HiddenKV[k].(string)
		return fl, nil
	case nil:
		return "", ErrValueWasNil
	// TODO: what about booleans?
	default:
		return "", fmt.Errorf("HiddenGetStringAlways: unrecognized type %v", v)
	}
}

// GetStringAlways is a safe, typed string accessor. It will force conversion away
// from float64, int64, uint64, string, and nil values. Nils and unrecognized values
// are marked as an error (nil values will return 0 and ErrValueWasNil)
func (o Object) GetStringAlways(k string) (string, error) {
	switch v := o.KV[k].(type) {
	case float64:
		fl := o.KV[k].(float64)
		return fmt.Sprintf("%f", fl), nil
	case int64:
		fl := o.KV[k].(int64)
		return fmt.Sprintf("%d", fl), nil
	case uint64:
		fl := o.KV[k].(uint64)
		return fmt.Sprintf("%d", fl), nil
	case string:
		fl := o.KV[k].(string)
		return fl, nil
	case nil:
		return "", ErrValueWasNil
	// TODO: what about booleans?
	default:
		return "", fmt.Errorf("GetStringAlways: unrecognized type %v", v)
	}
}

// GetFloatAlways is a safe, typed float64 accessor. It will force conversion away
// from float64, int64, uint64, string, and nil values. Nils and unrecognized values
// are marked as an error (nil values will return 0 and ErrValueWasNil)
func (o Object) GetFloatAlways(k string) (float64, error) {
	switch v := o.KV[k].(type) {
	case float64:
		fl := o.KV[k].(float64)
		return fl, nil
	case int64:
		fl := o.KV[k].(float64)
		return fl, nil
	case uint64:
		fl := o.KV[k].(uint64)
		return float64(fl), nil
	case string:
		fl := o.KV[k].(string)
		return strconv.ParseFloat(fl, 64)
	case nil:
		return 0, ErrValueWasNil
	// TODO: what about booleans?
	default:
		return 0, fmt.Errorf("GetFloatAlways: unrecognized type %v", v)
	}
}

// GetIntAlways is a safe, typed int64 accessor. It will force conversion away
// from float64, uint64, int64 and string values. Nils and unrecognized values are
// marked as an error (nil values will return 0 and ErrValueWasNil)
func (o Object) GetIntAlways(k string) (int64, error) {
	switch v := o.KV[k].(type) {
	case float64:
		fl := o.KV[k].(float64)
		return int64(fl), nil
	case int64:
		fl := o.KV[k].(int64)
		return fl, nil
	case uint64:
		fl := o.KV[k].(uint64)
		return int64(fl), nil
	case string:
		fl := o.KV[k].(string)
		return strconv.ParseInt(fl, 10, 64)
	case nil:
		return 0, ErrValueWasNil
	// TODO: what about booleans?
	default:
		return 0, fmt.Errorf("GetIntAlways: unrecognized type %v", v)
	}
}

// GetUintAlways is a safe, typed uint64 accessor. It will force conversion
// away from float64, int64, and string values. Nils and unrecognized values
// are marked as an error (nil values will return 0 and ErrValueWasNil)
func (o Object) GetUintAlways(k string) (uint64, error) {
	switch v := o.KV[k].(type) {
	case float64:
		fl := o.KV[k].(float64)
		return uint64(fl), nil
	case int64:
		fl := o.KV[k].(int64)
		return uint64(fl), nil
	case uint64:
		fl := o.KV[k].(uint64)
		return fl, nil
	case string:
		fl := o.KV[k].(string)
		return strconv.ParseUint(fl, 10, 64)
	case nil:
		return uint64(0), ErrValueWasNil
	// TODO: what about booleans?
	default:
		return 0, fmt.Errorf("GetUintAlways: unrecognized type %v", v)
	}
}

// Set is our typical setter. It attempts to track changes in records and the
// current state of whether an object appears to have been modified from what
// the database had (or should have).
func (o *Object) Set(k string, v interface{}) {
	oldVal := o.Get(k)

	if oldVal != nil {
		// Avoid redundant Set()s
		if oldVal == v {
			return
		}
		o.FieldChanged(k, oldVal)
	}
	if o.GetSaved() {
		o.SetSaved(false)
	}
	o.SetCore(k, v)
}

// FieldChanged records the previous value for something that is about to be
// set
func (o *Object) FieldChanged(k string, oldVal interface{}) {
	o.ChangedFields[k] = oldVal
}

// SetCore just mutates the internal object KV without any of the usual
// tracking that occurs when Set is called.
func (o *Object) SetCore(k string, v interface{}) {
	o.KV[k] = v
}

// ResetChangedFields can be used in conjunction with an ORM... For instance,
// once a Save() method is invoked
func (o *Object) ResetChangedFields() {
	o.ChangedFields = make(map[string]interface{})
}

// SetSaved is our setter for the 'object is saved' status field
func (o *Object) SetSaved(status bool) {
	o.saved = status
}

// GetSaved is our getter for the 'object is saved' bool field
func (o Object) GetSaved() bool {
	return o.saved
}
