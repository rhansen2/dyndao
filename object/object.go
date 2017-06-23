// Package object is an abstract data record which tracks state changes.
// The state change tracking can be useful when the values of primary keys
// need to change.
package object

// Object struct encapsulates our key-value pairs and a single-item per-key history
// of the previous value stored for a given key.
type Object struct {
	Type          string
	KV            map[string]interface{} `json:"KV"`
	ChangedFields map[string]interface{} `json:"ChangedFields"`
	Children      map[string]*Object
	saved         bool
}

// New is an empty constructor
func New() *Object {
	return &Object{KV: makeEmptyMap(), ChangedFields: makeEmptyMap(), Children: makeEmptyChildrenMap()}
}

func makeEmptyMap() map[string]interface{} {
	return make(map[string]interface{})
}

func makeEmptyChildrenMap() map[string]*Object {
	return make(map[string]*Object)
}

// TODO: ForEach method for the KV? ...

// Get is our accessor
func (o Object) Get(k string) interface{} {
	return o.KV[k]
}

// Set is our setter
func (o Object) Set(k string, v interface{}) {
	oldVal := o.Get(k)

	if oldVal != nil {
		// Avoid redundant Set()s
		if oldVal == v {
			return
		}
		o.FieldChanged(k, oldVal)
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
func (o Object) ResetChangedFields() {
	o.ChangedFields = make(map[string]interface{})
}

// SetSaved is our setter for the 'object is saved' status field
func (o Object) SetSaved(status bool) {
	o.saved = status
}

// GetSaved is our getter for the 'object is saved' bool field
func (o Object) GetSaved() bool {
	return o.saved
}
