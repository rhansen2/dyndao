package core

// TODO This package is bad. The functions imply that they will perform the
// action but they just return SQL.
// 

// Begin renders a SQL BEGIN statement for us
func Begin(name string) string {
	return "BEGIN " + name
}

// Commit returns a SQL Commit statement for us
func Commit() string {
	return "COMMIT"
}
