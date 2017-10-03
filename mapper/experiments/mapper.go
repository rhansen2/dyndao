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

func unmarshalToObjects(sch *schema.Schema, unm interface{}, rootLevel bool) (object.Array, error) {
	// top level is map
	m := unm.(map[string]interface{})
	var objAry object.Array
	for k, v := range m {

		switch v.(type) {
		case []interface{}:
			vi := v.([]interface{})
			objs, err := mapValToObjAry(sch, k, vi, false)
			if err != nil {
				return nil, err
			}
			objAry = append(objAry, objs...)
			continue
		case map[string]interface{}:
			vals := v.(map[string]interface{})

			obj := object.New(k)
			obj.KV = vals

			for childTableName := range sch.Tables {
				if obj.KV[childTableName] != nil {
					thingy := obj.KV[childTableName]
					interSlice := []interface{}{thingy}
					objs, err := mapValToObjAry(sch, k, interSlice, false)
					if err != nil {
						return nil, err
					}

					obj.Children[childTableName] = append(obj.Children[childTableName], objs...)
					delete(obj.KV, childTableName)
				}
			}

			objAry = append(objAry, obj)
		default:
			// TODO: how to note these to the caller?
			fmt.Println("Unrecognized: default: k=", k, "v=", v)
		}
	}

	return objAry, nil
}

// UnmarshalToObject accepts a schema configuration and a string of json data and
// returns an object array.
func UnmarshalToObject(sch *schema.Schema, json string) (*object.Object, error) {
	if json == "" {
		return nil, errors.New("ToObjectFromJSON: json parameter is empty")
	}

	var obj object.Object
	err := gjson.Unmarshal([]byte(json), &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

// ToObjectArrayFromJSON accepts a schema configuration and a string of json data and
// returns an object array.
func ToObjectArrayFromJSON(sch *schema.Schema, json string) (object.Array, error) {
	if json == "" {
		return nil, errors.New("ToObjectsFromJSON: json parameter is empty")
	}

	var unmarsh object.Array
	err := gjson.Unmarshal([]byte(json), &unmarsh)
	if err != nil {
		return nil, err
	}

	return unmarsh, nil
}

// ToObjectsFromJSON accepts a schema configuration and a string of json data and
// returns an object array.
func ToObjectsFromJSON(sch *schema.Schema, json string) (object.Array, error) {
	if json == "" {
		return nil, errors.New("ToObjectsFromJSON: json parameter is empty")
	}

	var unmarsh interface{}
	err := gjson.Unmarshal([]byte(json), &unmarsh)
	if err != nil {
		return nil, err
	}

	return unmarshalToObjects(sch, unmarsh, true)
}

func mapValToObjAry(sch *schema.Schema, objectType string, vals []interface{}, rootLevel bool) (object.Array, error) {
	var allObjs object.Array
	for _, val := range vals {
		switch t := val.(type) {
		case []interface{}:
			valAry := val.([]interface{})
			objAry, err := mapValToObjAry(sch, objectType, valAry, rootLevel)
			return objAry, err
		case map[string]interface{}:
			obj := object.New(objectType)
			m := val.(map[string]interface{})
			schTable := sch.Tables[objectType]

			for k, v := range m {
				if schTable.Children[k] == nil {
					obj.Set(k, v)
				} else {
					vAry := v.([]interface{})
					objAry, err := mapValToObjAry(sch, k, vAry, false)
					if err != nil {
						return nil, err
					}
					obj.Children[k] = append(obj.Children[k], objAry...)
				}
			}
			return object.NewArray(obj), nil
		default:
			// TODO:
			fmt.Println("[mapValToObjAry] Unrecognized: val ", val, " is type ", t)
		}
	}
	return allObjs, nil
}
