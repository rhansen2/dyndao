Installing DB2 on Ubuntu 16.04 was quite painful for me when I tried the Docker approach.

Instead, I recommend following this article:

https://globalengineer.wordpress.com/2017/02/27/installing-db2-on-debianubuntu/

and downloading DB2 directly from IBM.

You'll end up with the sqllib folder installed in /home/db2inst1 by default, so you can reference that
when installing github.com/alexbrainman/odbc

```code
#!/bin/bash

export DB2HOME=/home/db2inst1/sqllib
export CGO_LDFLAGS=-L$DB2HOME/lib
export CGO_CFLAGS=-I$DB2HOME/include

# now run go build . and go install . in your src/github.com/alexbrainman/odbc

```

Then to actually get the driver working:

1) Install unixODBC.
2) Run the script in odbc_ini to configure unixODBC for DB2

Now 'dyndao make test' if you have a dyndao alias setup with appropriate configuration.

--

To start db2 again after you've restarted your machine, sudo su - db2inst1, and run db2start.
