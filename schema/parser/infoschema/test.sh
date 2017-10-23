#!/bin/bash
DRIVER=mysql DSN=root:$SECRET_PW@/test?charset=utf8 go test -v
