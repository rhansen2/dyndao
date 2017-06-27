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
func ToJSONFromObject(sch *schema.Schema, obj *object.Object, rootJSON string, rootPath string, setRootPath bool) (string, error) {
	tbl := obj.Type
	table := sch.Tables[tbl]
	fieldsMap := table.Fields

	if rootJSON == "" {
		rootJSON = "{}"
	}
	if setRootPath {
		if rootPath == "" {
			rootPath = tbl
		} else {
			rootPath += "." + tbl
		}
		rootPath = rootPath + "."
	}
	// Populate object's key-value pairs into JSON
	for k, v := range obj.KV {
		field, ok := fieldsMap[k]
		if !ok {
			return "", errors.New("ToJSONFromObject: empty field for field " + k)
		}

		var err error
		// TODO: use table.JSONRoot or something instead of tbl here?
		rootJSON, err = sjson.Set(rootJSON, rootPath+field.Source, v)
		if err != nil {
			return "", err
		}
	}

	if obj.Children != nil && len(obj.Children) > 0 {
		for k, v := range obj.Children {
			for _, childObj := range v {
				childJSON, err := ToJSONFromObject(sch, childObj, "", "", false)
				if err != nil {
					return "", fmt.Errorf("ToObjectFromJSON: error %s with child [k:%v v:%v]", err.Error(), k, v)
				}
				// TODO: use i iterator variable somewhere here with SetRaw.. ? bench that?

				rootJSON, err = sjson.SetRaw(rootJSON, rootPath+k, childJSON)
				if err != nil {
					return "", fmt.Errorf("ToObjectFromJSON: error %s with child [k:%v v:%v]", err.Error(), k, v)
				}
			}
		}
	}
	return rootJSON, nil
}

func getPathPrefix(consumedPath string, tbl string) string {
	pathPrefix := ""
	if consumedPath != "" {
		pathPrefix = consumedPath
	}
	if pathPrefix == "" {
		pathPrefix = tbl
	} else {
		pathPrefix = pathPrefix + "." + tbl
	}
	return pathPrefix
}

// ToObjectFromJSON maps a JSON string into a DynDAO object
func ToObjectFromJSON(sch *schema.Schema, consumedPath string, tbl string, json string) (*object.Object, error) {
	obj := object.New(tbl)

	table := sch.Tables[tbl]
	fieldsMap := table.Fields

	keys := make([]string, len(fieldsMap))    // list of keys we're going to set
	sources := make([]string, len(fieldsMap)) // list of paths we'll retrieve them from

	pathPrefix := getPathPrefix(consumedPath, tbl)
	i := 0
	for k, field := range fieldsMap {
		if field.Source == "" {
			return nil, errors.New("ToObjectFromJSON: missing Source for field " + k)
		}
		keys[i] = k
		sources[i] = pathPrefix + "." + field.Source
		i++
	}
	values := gjson.GetMany(json, sources...)
	for i, v := range values {

		if v.Exists() {
			obj.Set(keys[i], v.Value())
		}
	}

	err := walkChildrenFromJSON(sch, table, obj, tbl, json)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func walkChildrenFromJSON(sch *schema.Schema, table *schema.Table, obj *object.Object, pathPrefix string, json string) error {
	if table.Children != nil && len(table.Children) > 0 {
		i := 0
		for k, v := range table.Children {
			jsonPath := k
			child, err := ToObjectFromJSON(sch, pathPrefix, jsonPath, json)
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

func ToObjectsFromJSON(sch *schema.Schema, json string) (object.Array, error) {
	if json == "" {
		return nil, errors.New("ToObjectsFromJSON: json parameter is empty")
	}
	parsed := gjson.Parse(json)
	objs := object.NewEmptyArray()
	if parsed.Type == gjson.JSON {
		m := parsed.Map()
		for k, v := range m {
			obj := object.New(k)

			if v.Type == gjson.JSON {
				parsed := gjson.Parse(v.Raw)

				parsed.ForEach(func(k, v gjson.Result) bool {
					ks := k.String()
					if v.Type == gjson.JSON {
						val := v.Value().([]interface{})
						objAry := mapValToObjAry(ks, val)
						obj.Children[ks] = objAry
					} else {
						obj.Set(ks, v.Value())
					}
					return true
				})
			}
			objs = append(objs, obj)
		}
	}

	return objs, nil
}

func mapValToObjAry(objectType string, vals []interface{}) object.Array {
	objs := make(object.Array, len(vals))
	for i, val := range vals {
		obj := object.New(objectType)
		m := val.(map[string]interface{})
		for k, v := range m {
			obj.Set(k, v)
		}
		objs[i] = obj
	}
	return objs
}
