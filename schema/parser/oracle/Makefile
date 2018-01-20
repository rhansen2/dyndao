# dyndao makefile, just for testing for now

# Just a test rule for now.
test:
	# do not use SYSTEM or SYS for the OWNER, use the actual schema-owning username.
	OWNER=username DRIVER=goracle ORACLE_SID=$(DBNAME) DSN=SYSTEM/oracle@//localhost:1521/XE go test -v 
