# dyndao
DYNamic Data Access Object (in Go)

JSON <-> object.Object <-> RDBMS

dyndao is a dynamic Golang ORM, drawing influence from the Active Record
pattern, Martin Fowler's Data Mapper pattern, and the author's experience with Perl's
database packages (DBI, Class::DBI, DBIx::Class, etc.)

Currently, basic support for the following databases is included:
* SQLite
* MySQL
* Microsoft SQL Server
* Oracle
* Postgres
* CockroachDB

If you encounter problems using dyndao with any of the currently supported
databases, please feel free to file a detailed issue.

The goal is not to support every feature of every database, but provide a
common set of data types and abstractions which allow a developer to deliver an
application that supports a myriad of underlying storage mechanisms.

*DATA TYPES, RELATIONAL JSON, ETC.*

If you have legacy requirements and specific data type requirements, dyndao may
or may not be the right solution for you. If you have more control over your
selection of data types, then it's more likely that dyndao could be a good fit.

The basic data types currently supported are: strings, integers, clobs / blobs,
floats, and timestamps.

*DYNAMIC SCHEMAS*

See github.com/rbastic/dyndao/schema for how dyndao supports dynamic schemas.

There are several options for the dyndao user. Declare a schema using Go code.
Edit one directly in JSON and have some code that unmarshals it. Or, if the
underlying database is well-supported, then dyndao's schema/parser sub-packages
can be utilized to dynamically load a schema at run-time.

*CODE LAYOUT*

```code
object - object.Object and object.Array

schema - schema definitions and supporting types

orm    - Bridge pattern-influenced package combining schema, sqlgen, and
	 object.

adapters - SQL statement generators for various database implementations

sqlgen - Specifies SQL Generator vtables

mapper - Custom JSON mapping layers (WIP)
```

*DISCLAIMER* 

This package is a work in progress. While much of the code is regularly tested
on production workloads, please use it at your own risk. Suggestions and
patches are welcome. Currently, I reserve the right to refactor the API at any
time.

If you find yourself using this package to do cool things, please let me know.
:-)

*THANKS*

The author would like to express his sincere thanks to Rob Hansen
(github.com/rhansen2), without whom this library surely would have suffered.

*LICENSE*

dyndao is released under the MIT license. See LICENSE file.

