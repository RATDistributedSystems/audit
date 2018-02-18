#!/bin/bash

cd ..
CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo --ldflags="-s" -o audit

# Build the image
docker build -t auditserver .

# Remove remnants
rm -f audit
