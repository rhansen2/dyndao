#!/bin/bash
 psql -h localhost -p 49153 -d docker -U $POSTGRES_USER --password
