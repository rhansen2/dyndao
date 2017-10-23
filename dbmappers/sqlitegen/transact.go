package sqlitegen

// Begin renders a SQL BEGIN statement for us
func (g Generator) Begin(name string) string {
	return "BEGIN " + name
}

// Commit returns a SQL Commit statement for us
func (g Generator) Commit() string {
	return "COMMIT"
}
