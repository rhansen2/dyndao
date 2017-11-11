The test code in here is a harness for the rest of the individual SQL Generators.

Each of the individual adapters (implementations of SQLGenerator) supplies
their own GetDB() and GetSQLGen() that is injected into the core test
framework's package namespace.

