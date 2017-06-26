// Package object is an abstract data record which tracks state changes.
// It's meant to make it easier to map key-value records into ORM / RDBMS
// structures.
// The state change tracking can be useful when the values of primary keys
// need to change. (Changing a foreign key on a table with a composite key, for example)
package object

// TODO: add object array container

// Object struct encapsulates our key-value pairs and a single-item per-key history
// of the previous value stored for a given key.

type ObjectArray []*Object

type Object struct {
	Type          string
	KV            map[string]interface{} `json:"KV"`
	ChangedFields map[string]interface{} `json:"ChangedFields"`
	Children      map[string]ObjectArray // TODO: make this value an array
	saved         bool
}

// New is an empty constructor
func New(objType string) *Object {
	return &Object{Type: objType, KV: makeEmptyMap(), ChangedFields: makeEmptyMap(), Children: makeEmptyChildrenMap(), saved: false}
}

func NewObjectArray(val *Object) ObjectArray {
	objAry := make(ObjectArray, 1)
	objAry[0] = val
	return objAry
}

func makeEmptyMap() map[string]interface{} {
	return make(map[string]interface{})
}

func makeEmptyChildrenMap() map[string]ObjectArray {
	return make(map[string]ObjectArray)
}

// Get is our accessor
func (o Object) Get(k string) interface{} {
	return o.KV[k]
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
