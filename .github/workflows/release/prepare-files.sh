#!/bin/bash

# version file
echo $1 > VERSION
echo $1 > ansible/VERSION
sed -i -E "s/(version: \").*(\")/\1$1\2/" ansible/galaxy.yml
