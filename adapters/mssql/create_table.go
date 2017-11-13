package mssql

import (
	"github.com/rbastic/dyndao/adapters/common"
	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

func RenderCreateColumn(sg *sg.SQLGenerator, f *schema.Column) string {
	return common.RenderCreateColumn(sg, f, "IDENTITY", mapType)
}

func mapType(s string) string {
	switch s {
	case "INTEGER":
		return "INT"
	case "BLOB":
		return "IMAGE"
	default:
		return s
	}
}
