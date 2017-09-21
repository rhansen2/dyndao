# dyndao
DYNamic Data Access Object (in Go)

JSON <-> object.Object <-> RDBMS

dyndao is a collection of Golang packages which comprise a JSON mapping layer,
a generic key-value object structure, and a relational database mapping layer
(ORM).

This package is a work in progress. While the package is regularly tested on
production workloads, please use it at your own risk. Suggestions and patches
are welcome. Please ensure any pull requests are go fmt'd and come up
(relatively) clean from gometalinter.

Currently, basic support for SQLite, MySQL, and Oracle are included.

*MOTIVATION*

Most ORMs perceive the database world as a static state of affairs. Go's
relatively static nature further complicates this situation.

Code generators alleviate some of the pain but in some situations, there is
still much left to be desired. dyndao is driven by the requirement for schemas
to be completely dynamic. This offers additionally flexibility at a cost to
performance.

See github.com/rbastic/dyndao/schema for how dyndao handles dynamic schemas.

Additionally, since most Go ORMs use code generators and structs, you are never
able to easily write an UPDATE that mentions a SQL function call:

```code
UPDATE fooTable SET ..., UPDATE_TIMESTAMP=NOW() WHERE fooTable_ID = 1;
```

NOW() is the timestamp function call (at least in MySQL, not sure of others).

Presently, the way dyndao is written, you could do something like:

```code
// 'o' is an instantiation of the ORM package
obj, err := o.RetrieveObject(ctx, tableString, pkValues)
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
	r, err := o.SaveObject(ctx, nil, obj)
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

JSON leaves much to be desired due to Go's type model. Mapping reliably between
JSON and an ORM is fraught with issues due to assumptions that Go makes with
regards to map[string]interface{}. Please take the opinions here with a grain of
salt, as one of the authors is originally a Perl programmer. As we all know,
Perl's core values are laziness, impatience, and hubris - hubris is a double-edged
sword.

Nonetheless, the authors of this package have attempted to manage these issues
where possible.

It may be necessary in some situations to store certain values as strings to
avoid the issues hinted at above. At least one such issue has resulted from
unsigned integer hash values that needed to be stored as strings (to avoid
being converted to Go's float64 type)

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

dyndao does not currently check contexts for expiration. You should do so
in your code before calling any of dyndao's ORM methods.

Example:
```code
select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
}

// now call method in question
```

*DISCLAIMER*

Please note that the current code layout and design is the result of deadlines,
hackathons, corporate feature requirements, and metric quantities of caffeine,
compounded over years of experience and multiple jaded developers. While humor
is intended, the author(s) can take no responsibility for damages that one may
incur by viewing or utilizing this source code.

*LICENSE*

MIT license. See LICENSE file.
