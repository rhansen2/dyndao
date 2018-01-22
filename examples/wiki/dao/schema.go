package dao

import (
	"github.com/rbastic/dyndao/schema"
	"sync"
)

const (
	SchInvoicesName = "SchInvoices"
	TblInvoicesName = "INVOICES"
)

var (
	activeSchema    *schema.Schema
	activeSchemaMtx sync.Mutex
)

func GetActiveSchema() *schema.Schema {
	return activeSchema
}

func SetActiveSchema(sch *schema.Schema) {
	activeSchemaMtx.Lock()
	defer activeSchemaMtx.Unlock()
	activeSchema = sch
}

func InvoiceSchema() *schema.Schema {
	sch := schema.DefaultSchema()
	sch.Name = SchInvoicesName

	tblInvoices := schema.DefaultTableWithName(TblInvoicesName)
	sch.Tables[TblInvoicesName] = tblInvoices

	return sch
}
