package sqlgen

// DropTable renders a SQL drop table statement for us
func (g Generator) DropTable() string {
	return "DROP TABLE"
}
