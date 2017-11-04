// Package common encapsulates common functionality that can be leveraged
// by all (or almost all) necessary database adapters. It differs from core
// in that common is still leveraged from the individual adapter code.
package common

import (
	"fmt"
	"strings"

	"github.com/rbastic/dyndao/schema"
	sg "github.com/rbastic/dyndao/sqlgen"
)

// TODO: AUTOINCREMENT attribute support.
// See http://www.sqlitetutorial.net/sqlite-autoincrement/

func RenderCreateColumn(sg *sg.SQLGenerator, f *schema.Column, identityStr string, mapTypeFn func(dbType string) string) string {
	dataType := strings.ToUpper(f.DBType)

	notNull := ""
	identity := ""
	unique := ""

	if f.IsIdentity {
		identity = identityStr
	}
	if f.AllowNull {
		notNull = "NULL"
	} else {
		notNull = "NOT NULL"
	}
	if mapTypeFn != nil {
		dataType = mapTypeFn(dataType)
	}
	if f.Length > 0 {
		dataType = fmt.Sprintf("%s(%d)", dataType, f.Length)
	}

	if f.IsUnique {
		unique = "UNIQUE"
	}

	return strings.Join([]string{f.Name, dataType, identity, notNull, unique}, " ")
}
