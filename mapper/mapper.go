// Package mapper is a DAO mapper for mapping between JSON, an generic object, and a configurable database schema
package mapper

import (
	"errors"
	"fmt"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// ToJSONFromObject encodes a DynDAO object
func ToJSONFromObject(sch *schema.Schema, obj *object.Object, jsonStr string) (string, error) {
	tbl := obj.Type
	table := sch.Tables[tbl]
	fieldsMap := table.Fields

	if jsonStr == "" {
		jsonStr = "{}"
	}
	for k, v := range obj.KV {
		field, ok := fieldsMap[k]
		if !ok {
			return "", errors.New("ToJSONFromObject: empty field for field " + k)
		}

		var err error
		jsonStr, err = sjson.Set(jsonStr, field.Source, v)
		if err != nil {
			return "", err
		}
	}

	if obj.Children != nil && len(obj.Children) > 0 {
		for k, v := range obj.Children {
			var err error
			jsonStr, err = ToJSONFromObject(sch, v, jsonStr)
			if err != nil {
				return "", fmt.Errorf("ToObjectFromJSON: error %s with child [k:%v v:%v]", err.Error(), k, v)
			}
			//			obj.Children[k] = child
		}
	}
	return jsonStr, nil
}

// ToObjectFromJSON maps a JSON string into a DynDAO object
func ToObjectFromJSON(sch *schema.Schema, tbl string, json string) (*object.Object, error) {
	obj := object.New(tbl)

	table := sch.Tables[tbl]
	fieldsMap := table.Fields

	keys := make([]string, len(fieldsMap))    // list of keys we're going to set
	sources := make([]string, len(fieldsMap)) // list of paths we'll retrieve them from

	i := 0
	for k, field := range fieldsMap {
		if field.Source == "" {
			return nil, errors.New("ToObjectFromJSON: missing Source for field " + k)
		}
		keys[i] = k
		sources[i] = field.Source
		i++
	}

	values := gjson.GetMany(json, sources...)

	for i, v := range values {
		if v.Exists() {
			obj.Set(keys[i], v.Value())
		} /* else {		return nil, errors.New("ToObjectFromJSON: missing value for field " + keys[i])}*/
	}

	err := walkChildrenFromJSON(sch, table, obj, json)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func walkChildrenFromJSON(sch *schema.Schema, table *schema.Table, obj *object.Object, json string) error {
	if table.Children != nil && len(table.Children) > 0 {
		for k, v := range table.Children {
			child, err := ToObjectFromJSON(sch, k, json)
			if err != nil {
				return fmt.Errorf("ToObjectFromJSON: error %s with child [k:%v v:%v]", err.Error(), k, v)
			}
			obj.Children[k] = child
		}
	}
	return nil
}
