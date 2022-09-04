#!/bin/bash

go version
go get -d
go generate ./...
go build -o ./bin/stackhead-cli .
