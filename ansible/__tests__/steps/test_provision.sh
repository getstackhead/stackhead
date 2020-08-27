#!/bin/bash
# IP address in environment "INPUT_IPADDRESS"
# Domain in environment "DOMAIN"
INVENTORY_PATH=${GITHUB_ACTION_PATH}/ansible/__tests__/inventory.yml

sed -e "s/\${ipaddress}/${INPUT_IPADDRESS}/" -e "s/\${webserver}/${INPUT_WEBSERVER}/" "${GITHUB_ACTION_PATH}/ansible/__tests__/inventory.dist.yml" > $INVENTORY_PATH

# Install dependencies
ansible-galaxy install -r "${GITHUB_ACTION_PATH}/ansible/requirements/requirements.yml" --force-with-deps

if [[ $INPUT_ROLENAME != '' ]]; then
  # Remove this role and set symlink
  rm -rf "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
  ln -s "${GITHUB_WORKSPACE}" "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
fi

# Provision server
TEST=1 ansible-playbook "${GITHUB_ACTION_PATH}/ansible/server-provision.yml" -i "${INVENTORY_PATH}" -vv
