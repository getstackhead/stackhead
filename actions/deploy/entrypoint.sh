#!/bin/sh -l

# Add IP to known_hosts
ssh-keyscan -v -T 30 ${2} >> ~/.ssh/known_hosts

# Install dependencies
/stackhead-cli init --version="${1}" -v

# Deploy project
/stackhead-cli project deploy "${GITHUB_WORKSPACE}/${3}" "${2}" -v

