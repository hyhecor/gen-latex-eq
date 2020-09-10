#!/bin/bash
## set version `git describe --tags $(git rev-list --tags --max-count=1)`
VERSION=`git describe --tags $(git rev-list --tags --max-count=1)`
## set build `git log -1 --pretty=format:%h`
BUILD=`git log -1 --pretty=format:%h`

## go install
go install -ldflags "-X main.version=${VERSION}@${BUILD}" 

## test 
gen-latex-eq -version