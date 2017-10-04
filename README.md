# dyndao
DYNamic Data Access Object (in Go)

JSON <-> object.Object <-> RDBMS

dyndao is a collection of Golang packages which serve as an
ActiveRecord-influenced ORM in Go.

*DISCLAIMER* 

This package is a work in progress. While the package is regularly tested on
production workloads, please use it at your own risk. Suggestions and patches
are welcome. I reserve the right to change the API at any time. Please ensure
any pull requests are go fmt'd and are (relatively) clean when run against
gometalinter.

Currently, basic support for SQLite, MySQL, and Oracle are included. If you
intend on working with MySQL, please be sure to check out the 'columntype'
branch for github.com/go-sql-driver/mysql

*MOTIVATION*

Most ORMs perceive the database world as a static state of affairs. Go's
relatively static nature makes writing an ORM a bit different than in more
dynamic languages, like Perl and Python.

Code generators alleviate some of the pain but in some situations, there is
still much left to be desired. dyndao is driven by the requirement for schemas
to be completely dynamic. This offers additionally flexibility at a cost to
performance.

See github.com/rbastic/dyndao/schema for how dyndao handles dynamic schemas.

Additionally, static typing fails when you want to write an UPDATE that
involves a SQL function call:

```code
UPDATE fooTable SET ..., UPDATE_TIMESTAMP=NOW() WHERE fooTable_ID = 1;
```

NOW() is the current timestamp function call (at least in MySQL, not sure of
others). There is an issue open currently to abstract this further.

Presently, the way dyndao is written, you could do something like:

```code
myORM = getORM() // you'll have to write this to instantiate a db connection
obj, err := myORM.RetrieveObject(ctx, tableString, pkValues)
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
details). Almost everything is an object in dyndao.

*MISC.*

JSON leaves much to be desired due to Go's type model.  It may be necessary in
some situations to store certain values as strings to avoid the issues hinted
at above. At least one such issue has resulted from unsigned integer hash
values that needed to be stored as strings within object.Object (to avoid being
converted to Go's float64 type), only to later be mapped to number fields in
the relevant database tables.

*CODE LAYOUT*

```code
mapper - JSON mapping layer (WIP)

object - object.Object and object.Array

orm    - ORM class, combines object, schema, and sqlgen

schema - dynamic schema packages: declare your schema using these, or write and
         share a schema parser / generator with us for your database of choice!

sqlgen - code generators for various database implementations
```

*CONTEXT CHECKING*

dyndao no longer requires you to check your own contexts.

*THANKS*

The author would like to express his sincere thanks to Rob Hansen (github.com/rhansen2),
without whom the Oracle support in this library would not be possible.

*LICENSE*

MIT license. See LICENSE file.

