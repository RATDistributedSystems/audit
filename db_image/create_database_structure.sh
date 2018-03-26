#!/bin/bash
(sleep 65 && cqlsh -f create_audit_structure.cql) &
cd / && docker-entrypoint.sh