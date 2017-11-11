package main

import (
	_ "bitbucket.org/phiggins/db2cli"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	connStr = flag.String("conn", "", "connection string to use")
	repeat  = flag.Uint("repeat", 1, "number of times to repeat query")
)

func usage() {
	fmt.Fprintf(os.Stderr, `usage: %s [options]

%s connects to DB2 and executes a simple SQL statement a configurable
number of times.

Here is a sample connection string:

DATABASE=MYDBNAME; HOSTNAME=localhost; PORT=60000; PROTOCOL=TCPIP; UID=username; PWD=password;
`, os.Args[0], os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func execQuery(st *sql.Stmt) error {
	rows, err := st.Query()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var t time.Time
		err = rows.Scan(&t)
		if err != nil {
			return err
		}
		fmt.Printf("Time: %v\n", t)
	}
	return rows.Err()
}

func dbOperations() error {
	db, err := sql.Open("db2-cli", *connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	// Attention: If you have to go through DB2-Connect you have to terminate SQL-statements with ';'
	st, err := db.Prepare("select current timestamp from sysibm.sysdummy1;")
	if err != nil {
		return err
	}
	defer st.Close()

	for i := 0; i < int(*repeat); i++ {
		err = execQuery(st)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if *connStr == "" {
		fmt.Fprintln(os.Stderr, "-conn is required")
		flag.Usage()
	}

	if err := dbOperations(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
