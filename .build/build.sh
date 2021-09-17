#!/bin/bash

go get -d
pkger
go generate ./...
go build -o ./bin/stackhead-cli .
