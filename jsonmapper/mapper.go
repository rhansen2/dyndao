// Package jsonmapper is a DAO mapper for mapping between JSON, an generic object, and a configurable database schema
package jsonmapper

import (
	"errors"

	"github.com/rbastic/dyndao/object"
	"github.com/tidwall/gjson"
)

// ToObjectArrayFromJSON accepts a string of json data and returns an object
// array.
func ToObjectArrayFromJSON(json string) (object.Array, error) {
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
