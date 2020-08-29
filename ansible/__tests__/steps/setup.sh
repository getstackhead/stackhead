#!/bin/bash
# IP address in environment "INPUT_IPADDRESS"
INVENTORY_PATH=${GITHUB_ACTION_PATH}/ansible/__tests__/inventory.yml

sed -e "s/\${ipaddress}/${INPUT_IPADDRESS}/" -e "s/\${webserver}/${INPUT_WEBSERVER}/" -e "s/\${container}/${INPUT_CONTAINER}/" "${GITHUB_ACTION_PATH}/ansible/__tests__/inventory.dist.yml" > "${INVENTORY_PATH}"

# Provision server
TEST=1 ansible-playbook "${GITHUB_ACTION_PATH}/ansible/playbooks/server-provision.yml" -i "${INVENTORY_PATH}" -vv
