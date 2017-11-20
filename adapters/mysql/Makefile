# dyndao makefile, just for testing for now

# Just a test rule for now.
test:
	MYSQL_DSN=root:$(MYSQL_PASS)@/test?charset=utf8 go test -v .

cover:
	MYSQL_DSN=root:$(MYSQL_PASS)@/test?charset=utf8 go test -cover

htmlcover:
	MYSQL_DSN=root:$(MYSQL_PASS)@/test?charset=utf8 go test -coverprofile=coverage.out ; go tool cover -html=coverage.out

bench:
	MYSQL_DSN=root:$(MYSQL_PASS)@/test?charset=utf8 go test -bench=.

race:
	MYSQL_DSN=root:$(MYSQL_PASS)@/test?charset=utf8 go test -race
