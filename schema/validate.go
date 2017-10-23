package schema

import (
	"fmt"
)

var (
	stdErrorFmt = "dyndao/schema/Validate: schema.Table named '%s' has error %s"
)

func errorHelper(tbl *Table, msg string) error {
	return fmt.Errorf(stdErrorFmt, tbl.Name, msg)
}

// Validate is a basic schema validator. It ensures that each table inside the
// schema has a name, some Fields, and EssentialFields is set. Any other database
// requirements are not yet considered.
func Validate(sch *Schema) error {
	for _, tbl := range sch.Tables {
		if tbl.Name == "" {
			return errorHelper(tbl, "empty Name property")
		}

		if tbl.EssentialFields == nil {
			return errorHelper(tbl, "EssentialFields is nil")
		}

		if len(tbl.EssentialFields) == 0 {
			return errorHelper(tbl, "EssentialFields is of length zero")
		}

		// TODO: What other requirements do we have for defining a valid
		// schema?
	}

	return nil
}
