#!/bin/bash
ORACLE_DSN=system/oracle@//localhost:1521/xe.oracle.docker go test -v .
