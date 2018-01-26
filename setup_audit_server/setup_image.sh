#!/bin/bash

cd ..
docker run --rm -it -v "$GOPATH":/gopath -v "$(pwd)":/app -e "GOPATH=/gopath" -w /app golang:1.9 sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o audit'

# Build the image
docker build -t audit_server .

# Remove remnants
rm -f audit