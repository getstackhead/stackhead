#!/bin/bash
# This test destroys a deployed project on the target server
# IMPORTANT: This must run after test_deploy.sh!

CI=1 ${INPUT_CLI_BIN_PATH} project destroy "${GITHUB_ACTION_PATH}/container.stackhead.yml" "${INPUT_IPADDRESS}" -c "/tmp/.stackhead-cli.yml" -v
