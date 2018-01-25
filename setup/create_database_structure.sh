#!/bin/bash
cd /
cd - && sleep 40 && cqlsh -f create_database_structure.cql &
docker-entrypoint.sh

