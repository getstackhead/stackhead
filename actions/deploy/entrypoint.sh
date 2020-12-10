#!/bin/sh -l

mkdir -p ~/.ssh
echo "${1}" > ~/.ssh/id_rsa
chmod 0600 ~/.ssh/id_rsa
eval `ssh-agent -s`
ssh-add ~/.ssh/id_rsa

# Install dependencies
/bin/stackhead-cli init --version="${4}" -v

# Deploy project
echo "Deploying ${GITHUB_WORKSPACE}/${3} to server ${2}"
/bin/stackhead-cli project deploy "${GITHUB_WORKSPACE}/${3}" "${2}" -v

