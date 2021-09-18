#!/bin/bash
# This test deploys a container project onto the target server
# IMPORTANT: This must run after test_provision.sh!

sed -e "s/\${domain}/${INPUT_DOMAIN}/" -e "s/\${domain2}/${INPUT_DOMAIN2}/" "${GITHUB_ACTION_PATH}/project-definition.dist.yml" >"${GITHUB_ACTION_PATH}/container.stackhead.yml"

CI=1 ${INPUT_CLI_BIN_PATH} project deploy "${GITHUB_ACTION_PATH}/container.stackhead.yml" "${INPUT_IPADDRESS}" -c "/tmp/.stackhead-cli.yml" -v
