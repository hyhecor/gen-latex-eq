#!/bin/bash
## set version `git describe --tag`
VERSION=`git describe --tags`
## set build `git log -1 --pretty=format:%h`
BUILD=`git log -1 --pretty=format:%h`

## go build
go build -ldflags "-X main.version=${VERSION} -X main.build=${BUILD}" 

## test version
./gen-latex-eq -version
## test function
./gen-latex-eq < test.list