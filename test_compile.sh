#!/bin/bash

set -e
pwd
go get github.com/BubbaJoe/spvwallet-cash
go get github.com/mattn/go-sqlite3
go test -coverprofile=bitcoincash.cover.out ./
echo "mode: set" > coverage.out && cat *.cover.out | grep -v mode: | sort -r | \
awk '{if($1 != last) {print $0;last=$1}}' >> coverage.out
rm -rf *.cover.out
