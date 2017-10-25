# dyndao
DYNamic Data Access Object (in Go)

JSON <-> object.Object <-> RDBMS

dyndao is a dynamic Golang ORM, drawing influence from the Active Record
pattern, Martin Fowler's Data Mapper pattern, and the author's experience with Perl's
database packages (DBI, Class::DBI, DBIx::Class, etc.)

Currently, basic support for SQLite, MySQL, and Oracle are included. If you
intend on working with MySQL, please be sure to check out (as in, `git
checkout`) the 'columntype' branch for github.com/go-sql-driver/mysql.  Support
for additional database drivers is planned.

dyndao started out due to the fact that Go lacks a 'very dynamic' ORM,
and the author had some requirements for dynamic schemas and a vision for
multiple database support.

Despite the additional cost of supporting a complex 'object' type in Go,
this package is presently suitable for the author's needs.

See github.com/rbastic/dyndao/schema for how dyndao handles dynamic schemas.

One nice thing that dyndao supports is a type named object.SQLValue,
which lets you explicitly store values that can be rendered as
SQL function calls (and other unquoted necessities), without making the
mistake of constructing them as binding parameters.

I specifically mention this as an example because I have not yet seen
an ORM support this particular feature in Go, and dyndao's design
lends itself to implementation being simple. Here is an example:

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
	r, err := myORM.SaveObject(ctx, nil, obj)
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

*MISC.*

JSON leaves much to be desired due to Go's type model.  It may be necessary in
some situations to store certain values as strings to avoid the issues hinted
at above. At least one such issue has resulted from unsigned integer hash
values that needed to be stored as strings within object.Object (to avoid being
converted to Go's float64 type), only to later be mapped to number fields in
the relevant database tables.

*CODE LAYOUT*

```code
object - object.Object and object.Array

schema - schema definitions and supporting types

orm    - Bridge pattern-influenced package combining schema, sqlgen, and
	 object.

sqlgen - SQL statement generators for various database implementations

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
(github.com/rhansen2), for without him this library surely would have suffered.

*LICENSE*

dyndao is released under the MIT license. See LICENSE file.

