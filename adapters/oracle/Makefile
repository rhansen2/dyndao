# dyndao makefile, just for testing for now

# Just a test rule for now.

#USER=SYSTEM
#PASS=oracle

USER=gdl_join_w2
PASS=gdl_join_w2

test:
	ORACLE_SID=GDL_CURRENT4 ORACLE_DSN=$(USER)/$(PASS)@//localhost:1521/XE go test -v 

cover:
	ORACLE_SID=GDL_CURRENT4 ORACLE_DSN=$(USER)/$(PASS)@//localhost:1521/XE go test -cover

bench:
	ORACLE_SID=GDL_CURRENT4 ORACLE_DSN=$(USER)/$(PASS)@//localhost:1521/XE go test -bench=.

race:
	ORACLE_SID=GDL_CURRENT4 ORACLE_DSN=$(USER)/$(PASS)@//localhost:1521/XE go test -race
