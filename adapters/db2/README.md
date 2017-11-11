Installing DB2 on Ubuntu 16.04 was quite painful for me when I tried the Docker approach.

Instead, I recommend following this article:

https://globalengineer.wordpress.com/2017/02/27/installing-db2-on-debianubuntu/

and downloading DB2 directly from IBM.

You'll end up with the sqllib folder installed in /home/db2inst1 by default, so you can reference that
when installing bitbucket.org/phiggins/db2cli

Here's an example of the file that I used to get db2cli to build.

```code
#!/bin/bash

DB2HOME=/home/db2inst1/sqllib
export CGO_LDFLAGS=-L$DB2HOME/lib
export CGO_CFLAGS=-I$DB2HOME/include

# now run go build . in your src/bitbucket.org/phiggins/db2cli

```
