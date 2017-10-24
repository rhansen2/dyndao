# dyndao makefile, just for testing for now

# Just a test rule for now.
test:
	DRIVER=mysql DSN=root:$(SECRET_PW)@/test?charset=utf8 go test -v

cover:
	DRIVER=mysql DSN=root:$(SECRET_PW)@/test?charset=utf8 go test -cover
