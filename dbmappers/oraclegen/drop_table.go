package oraclegen

// DropTable renders a SQL drop table statement for us
func DropTable(name string) string {
	return "DROP TABLE " + name
}
