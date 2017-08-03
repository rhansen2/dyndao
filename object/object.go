// Package object is an abstract data record which tracks state changes.
// It's meant to make it easier to map key-value records into ORM / RDBMS
// structures.
// The state change tracking can be useful when the values of primary keys
// need to change. (Changing a foreign key on a table with a composite key, for example)
package object

import (
	"fmt"
)

// Array is our object array container to facilitate a couple of instances
// where slices are needed.
type Array []*Object

// Object struct encapsulates our key-value pairs and a single-item per-key history
// of the previous value stored for a given key.
type Object struct {
	Type          string
	KVOnlyUpper   bool
	KV            map[string]interface{} `json:"KV"`
	ChangedFields map[string]interface{} `json:"ChangedFields"`
	Children      map[string]Array
	saved         bool
}

// New is an empty constructor
func New(objType string) *Object {
	return &Object{Type: objType, KV: makeEmptyMap(), ChangedFields: makeEmptyMap(), Children: makeEmptyChildrenMap(), saved: false}
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

func makeEmptyMap() map[string]interface{} {
	return make(map[string]interface{})
}

func makeEmptyChildrenMap() map[string]Array {
	return make(map[string]Array)
}

// Get is our accessor
func (o Object) Get(k string) interface{} {
	return o.KV[k]
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

// GetFloat is a safe, typed int64 accessor
func (o Object) GetFloat(k string) (float64, bool) {
	v, ok := o.KV[k].(float64)
	return v, ok
}

// GetIntAlways is a safe, typed int64 accessor. It will force conversion away
// from float64 and uint64 values.
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
	default:
		return 0, fmt.Errorf("GetIntAlways: unrecognized type %v", v)
	}
}

// GetUintAlways is a safe, typed uint64 accessor. It will force conversion away
// from float64 and int64 values.
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
	default:
		return 0, fmt.Errorf("GetIntAlways: unrecognized type %v", v)
	}
}

// Set is our setter
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
	o.rawSet(k, v)
}

// FieldChanged records the previous value for something that is about to be set
func (o Object) FieldChanged(k string, oldVal interface{}) {
	o.ChangedFields[k] = oldVal
}

func (o Object) rawSet(k string, v interface{}) {
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
