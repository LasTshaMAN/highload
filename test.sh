#!/usr/bin/env bash

# delete old '*gen.go' files
find . -name '*gen.go' | xargs rm -rf -v

# generate new '*gen.go' files
egrep -l -R --include '*.go' "(go:generate)" . | sed 's/^/go generate /' | sh -v

# run all the tests in the project
find . -name '*_test.go' | xargs -n1 dirname | xargs go test -race
