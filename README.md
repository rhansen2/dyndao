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

If you intend on working with MySQL, please be sure to check out
(as in, `git checkout`) the 'columntype' branch for
github.com/go-sql-driver/mysql.  Support for additional database drivers is
planned. If you encounter problems using dyndao with any of the currently
supported databases, please feel free to file a detailed issue.

dyndao started out due to the fact that Go lacks a 'very dynamic' ORM, and the
author had some requirements for dynamic schemas and a vision for multiple
database support. The goal is not to support every feature of every database,
but provide a common set of data types and abstractions which allow a developer
to deliver an application that supports a myriad of underlying storage
mechanisms.

*DATA TYPES, RELATIONAL JSON, ETC.*

If you have legacy requirements and specific data type requirements, dyndao may
or may not be the right solution for you. If you have more control over your
selection of data types, then it's more likely that dyndao could be a good fit.

The basic data types currently supported are: strings, integers, clobs / blobs,
floats, and timestamps.

The reasoning here is that dyndao is geared towards systems which utilize a
'relational JSON' approach - the trend of only defining certain columns to
exist on a table, typically primary keys and timestamps, and then adding a JSON
document to support an arbitrary number of additional data attributes.

Thus, not all additions or removal of data attributes will require schema
changes, because the JSON document column can be leveraged to store them. This
provides some of the benefits of NoSQL in conjunction with the relational and
ACID benefits of SQL. 

Since many relational databases are adding support for indexing JSON documents
and directly querying them in a WHERE clause, likely as a response to competing
with NoSQL solutions, this idea is believed to be sound.

Additionally, many relational databases support data types which offer
inconsistent functionality despite having similar names. The possible length of
one column data type for a given database may be completely different from
another (TODO: provide examples).

*DYNAMIC SCHEMAS*

See github.com/rbastic/dyndao/schema for how dyndao supports dynamic schemas.

There are several options for the dyndao user. Declare a schema using Go code.
Edit one directly in JSON and have some code that unmarshals it. Or, if the
underlying database being used is well-supported, then dyndao's schema/parser
sub-packages can be utilized to dynamically load a schema at run-time.

*MISC*

An additional thing that dyndao supports is a type named object.SQLValue,
which lets you explicitly store values that can be rendered as
SQL function calls (and other unquoted necessities), without making the
mistake of constructing them as binding parameters.

NULL values mapped from the database will also be internally mapped using the
object.NewNULLValue() method. This is in contrast to how Go's sql.NullString,
etc., work. Those types are still utilized underneath the hood where necessary,
but to provide a common abstraction, the type object.SQLValue is used, with the
internal SQLValue string being set with a value of "NULL".

I specifically mention this as an example because I have not yet seen
an ORM support this particular feature in Go, and dyndao's design
lends itself to using the object.SQLValue for other things, which are not yet
covered here.  Here is an example of the initial use case:

```code
UPDATE fooTable SET ..., UPDATE_TIMESTAMP=NOW() WHERE fooTable_ID = 1;
```

NOW() is the current timestamp function call (at least in MySQL).

Presently, the way dyndao is written, you could do something like:

```code
myORM = getORM() // see tests for example
obj, err := myORM.Retrieve(ctx, tableString, pkValues)
if err != nil {
	panic(err)
}
// Lack of result is not considered an error in dyndao - "nil, nil" would be
// returned
if obj == nil {
	panic("no object available to update")
} else {
	// obj is a dyndao/object, see github.com/rbastic/dyndao/object
	obj.Set("UPDATE_TIMESTAMP", object.NewSQLValue("NOW()"))
	// nil means 'transactionless save', otherwise you can pass
	// a *sql.Tx
	r, err := myORM.Save(ctx, nil, obj)
	if err != nil {
		panic(err)
	}
	// ...
}
```

So, instead of representing rows as structs like other ORMs, you represent them
as pointers to object.Object (see github.com/rbastic/dyndao/object for
details).

Note, one should not have to hard-code the NOW() function, because if the
underlying database changes to say, Oracle, it will not work. Some work is
underway to explore how better cross-platform support for date- and timestamp-
handling can be implemented.

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

*CONTEXT CHECKING*

dyndao no longer requires you to check your own contexts.

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

