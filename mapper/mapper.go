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
func ToJSONFromObject(sch *schema.Schema, obj *object.Object, rootJSON string, rootPath string) (string, error) {
	tbl := obj.Type
	table := sch.Tables[tbl]
	fieldsMap := table.Fields

	if rootJSON == "" {
		rootJSON = "{}"
	}
	if rootPath == "" {
		rootPath = tbl
	} else {
		rootPath += "." + tbl
	}

	for k, v := range obj.KV {
		field, ok := fieldsMap[k]
		if !ok {
			return "", errors.New("ToJSONFromObject: empty field for field " + k)
		}

		var err error
		// TODO: use table.JSONRoot or something instead of tbl here?
		rootJSON, err = sjson.Set(rootJSON, field.Source, v)
		if err != nil {
			return "", err
		}
	}

	if obj.Children != nil && len(obj.Children) > 0 {
		for k, v := range obj.Children {
			for _, childObj := range v {
				childJSON, err := ToJSONFromObject(sch, childObj, "{}", rootPath)
				if err != nil {
					return "", fmt.Errorf("ToObjectFromJSON: error %s with child [k:%v v:%v]", err.Error(), k, v)
				}
				// TODO: use i iterator variable somewhere here with SetRaw..
				rootJSON, err = sjson.SetRaw(rootJSON, rootPath+"."+k, childJSON)
				if err != nil {
					return "", fmt.Errorf("ToObjectFromJSON: error %s with child [k:%v v:%v]", err.Error(), k, v)
				}
			}
		}
	}
	return rootJSON, nil
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
		}
	}

	err := walkChildrenFromJSON(sch, table, obj, json)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func walkChildrenFromJSON(sch *schema.Schema, table *schema.Table, obj *object.Object, json string) error {
	if table.Children != nil && len(table.Children) > 0 {
		i := 0
		for k, v := range table.Children {
			child, err := ToObjectFromJSON(sch, k, json)
			if err != nil {
				return fmt.Errorf("ToObjectFromJSON: error %s with child [k:%v v:%v]", err.Error(), k, v)
			}
			if obj.Children[k] == nil {
				obj.Children[k] = make(object.Array, len(table.Children))
			}
			obj.Children[k][i] = child
			i++
		}
	}
	return nil
}
