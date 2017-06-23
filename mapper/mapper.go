// Package mapper is a DAO mapper for mapping between an object, a schema, and a database generator (sqlgen)
package mapper

import (
	"errors"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/tidwall/gjson"
)

func toObjectFromJSON(sch *schema.Schema, tbl string, json string) (*object.Object, error) {
	obj := object.New()
	obj.Type = tbl // NOTE: we use the database table as the object 'type'

	fieldsMap := sch.Tables[tbl].Fields

	keys := make([]string, len(fieldsMap))    // list of keys we're going to set
	sources := make([]string, len(fieldsMap)) // list of paths we'll retrieve them from
	i := 0
	for k, field := range fieldsMap {
		if field.Source == "" {
			return nil, errors.New("toObjectFromJSON: empty Source for field " + k)
		}

		keys[i] = k
		sources[i] = field.Source
		i++
	}

	values := gjson.GetMany(json, sources...)

	for i, v := range values {
		if v.Exists() {
			obj.Set(keys[i], v)
		} else {
			return nil, errors.New("ToObjectFromJSON: empty value for field " + keys[i])
		}
	}

	// TODO: Implement Children piece

	return obj, nil
}
