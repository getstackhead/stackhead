#!/bin/bash

# version file
echo $1 > VERSION
echo $1 > ansible/VERSION

# update schemas in collection
(cd schemas && find . -name "*.json" -exec cp --parents -R '{}' "../ansible/schemas/" ';')
