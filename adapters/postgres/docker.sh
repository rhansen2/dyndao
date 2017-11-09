#!/bin/bash
POSTGRES_USER=postgres docker run --name magic_school_bus --user postgres -e POSTGRES_PASSWORD=dyndaoPassword -d postgres

