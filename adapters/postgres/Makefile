# dyndao makefile, just for testing for now

# Just a test rule for now.
test:
	POSTGRES_DSN='user=$(POSTGRES_USER) password=$(POSTGRES_PASS) dbname=$(POSTGRES_DB) sslmode=disable' go test -v 

cover:
	POSTGRES_DSN='user=$(POSTGRES_USER) password=$(POSTGRES_PASS) dbname=$(POSTGRES_DB) sslmode=disable' go test -cover
