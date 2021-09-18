#!/bin/bash
# IP address in environment "INPUT_IPADDRESS"

CI=1 ${INPUT_CLI_BIN_PATH} setup "${INPUT_IPADDRESS}" -c "/tmp/.stackhead-cli.yml" -v
