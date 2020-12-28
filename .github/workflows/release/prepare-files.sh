#!/bin/bash

# version file
echo "v$1" > VERSION
echo "v$1" > ansible/VERSION

# update schemas in collection
(cd schemas && find . -name "*.json" -exec cp --parents -R '{}' "../ansible/schemas/" ';')
