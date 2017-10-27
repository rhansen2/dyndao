# dyndao makefile, just for testing for now

# Just a test rule for now.
test:
	# This DSN broke things.
	#MSSQL_DSN='sqlserver://$(MSSQL_USER):$(MSSQL_PASS)@$(MSSQL_HOST):1433/instance?database=test' go test -cover
	MSSQL_DSN='sqlserver://$(MSSQL_USER):$(MSSQL_PASS)@$(MSSQL_HOST):1433' go test -v 

cover:
	
	MSSQL_DSN='sqlserver://$(MSSQL_USER):$(MSSQL_PASS)@$(MSSQL_HOST):1433' go test -cover
