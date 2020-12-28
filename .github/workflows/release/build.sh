#!/bin/bash

set -euo pipefail

pluginDir=".semrel/$(go env GOOS)_$(go env GOARCH)/hooks-logger/1.0.0/"
[[ ! -d "$pluginDir" ]] && {
  echo "creating $pluginDir"
  mkdir -p $pluginDir
}

go build -o $pluginDir/logger "$(pwd)/$(dirname $0)/main.go"
