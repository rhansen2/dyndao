package sqlitegen

// DropTable renders a SQL drop table statement for us
func (g Generator) DropTable(name string) string {
	return "DROP TABLE " + name
}
