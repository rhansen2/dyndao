package sqlite

// Begin renders a SQL BEGIN statement for us
func Begin(name string) string {
	return "BEGIN " + name
}

// Commit returns a SQL Commit statement for us
func Commit() string {
	return "COMMIT"
}
