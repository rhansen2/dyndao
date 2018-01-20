# dyndao makefile, just for testing for now

# Just a test rule for now.
test:
	ROACH_DSN='postgresql://$(ROACH_USER)@$(ROACH_HOST):$(ROACH_PORT)/$(ROACH_DB)?$(ROACH_PARAMS)' go test -v 

cover:
	ROACH_DSN='postgresql://$(ROACH_USER)@$(ROACH_HOST):$(ROACH_PORT)/$(ROACH_DB)?$(ROACH_PARAMS)' go test -cover

bench:
	ROACH_DSN='postgresql://$(ROACH_USER)@$(ROACH_HOST):$(ROACH_PORT)/$(ROACH_DB)?$(ROACH_PARAMS)' go test -bench=.

race:
	ROACH_DSN='postgresql://$(ROACH_USER)@$(ROACH_HOST):$(ROACH_PORT)/$(ROACH_DB)?$(ROACH_PARAMS)' go test -race
