#!/bin/bash
# IP address in environment "INPUT_IPADDRESS"

if [[ $INPUT_CLI != '' ]]; then
  ${INPUT_CLI_BIN_PATH} setup "${INPUT_IPADDRESS}" -c "/tmp/.stackhead-cli.yml"
else
  INVENTORY_PATH=${GITHUB_ACTION_PATH}/ansible/__tests__/inventory.yml
  sed -e "s/\${ipaddress}/${INPUT_IPADDRESS}/" -e "s/\${webserver}/${INPUT_WEBSERVER}/" -e "s/\${container}/${INPUT_CONTAINER}/" "${GITHUB_ACTION_PATH}/ansible/__tests__/inventory.dist.yml" > "${INVENTORY_PATH}"

  # Provision server
  TEST=1 ansible-playbook "${GITHUB_ACTION_PATH}/ansible/playbooks/server-provision.yml" -i "${INVENTORY_PATH}" -vv
fi
