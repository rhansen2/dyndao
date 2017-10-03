// Package json is a DAO mapper for mapping between JSON, an generic object, and a configurable database schema
package json

import (
	"errors"
	"fmt"

	"github.com/rbastic/dyndao/object"
	"github.com/rbastic/dyndao/schema"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var (
	OpCodeDTM
)

type JSONMappings []JSONMapping

type JSONMapping struct {
	OpCode int
	
}

type JSONMapper struct {
	Mappings JSONMappings
}


