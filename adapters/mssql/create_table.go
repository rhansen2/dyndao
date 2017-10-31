package mssql

import (
	"github.com/rbastic/dyndao/schema"
	"github.com/rbastic/dyndao/adapters/common"
	sg "github.com/rbastic/dyndao/sqlgen"
)

func RenderCreateColumn(sg *sg.SQLGenerator, f *schema.Column) string {
	return common.RenderCreateColumn(sg, f, "IDENTITY", mapType)
}

func mapType(s string) string {
	// Map 'integer' to 'number' for now for Oracle
	if s == "INTEGER" {
		return "INT"
	}
	if s == "TEXT" {
		return "TEXT"
	}
	if s == "BLOB" {
		return "IMAGE"
	}
	return s
}
