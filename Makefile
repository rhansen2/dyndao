# Set an alias in your .bashrc similar to this:
#
# alias dyndao='MYSQL_PASS=dyndaoPassword MSSQL_USER=dyndao MSSQL_PASS=dyndaoPassword MSSQL_HOST=127.0.0.1 POSTGRES_USER=docker POSTGRES_PASS=docker POSTGRES_DB=postgres DB2_USER=db2inst1 DB2_PASS=dyndaoPassword DB2_HOST=127.0.0.1 DB2_PORT=50000'
#
# Then source ~/.bashrc and run 'dyndao make test'
#

# Adapters not listed below are not yet 'finished'
test:
	cd adapters/sqlite; make test
	cd adapters/mysql; make test
	cd adapters/oracle; make test
	cd adapters/mssql; make test
	cd adapters/postgres; make test

race:
	cd adapters/sqlite; TEST_RACE=1 make race
	cd adapters/mysql; TEST_RACE=1 make race
	cd adapters/oracle; TEST_RACE=1 make race
	cd adapters/mssql; TEST_RACE=1 make race
	cd adapters/postgres; TEST_RACE=1 make race

cover:
	cd adapters/sqlite;  make cover
	cd adapters/mysql;  make cover
	cd adapters/oracle;  make cover
	cd adapters/mssql;  make cover
	cd adapters/postgres;  make cover
