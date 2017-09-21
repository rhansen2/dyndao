#!/bin/bash
DSN=root:$SECRET_PW@/test?charset=utf8 go test -v .
