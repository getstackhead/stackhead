#!/bin/bash
# IP address in environment "INPUT_IPADDRESS"
# Domain in environment "DOMAIN"
INVENTORY_PATH=${ACTION_PATH}/tests/inventory.yml

sed -e "s/\${ipaddress}/${INPUT_IPADDRESS}/" -e "s/\${webserver}/${INPUT_WEBSERVER}/" ${ACTION_PATH}/tests/inventory.dist.yml > $INVENTORY_PATH

echo "${GITHUB_ACTION_PATH} - ${ACTION_PATH}"

# Install dependencies
ansible-galaxy install -r "${ACTION_PATH}/ansible/requirements/requirements.yml" --force-with-deps

if [[ $INPUT_ROLENAME != '' ]]; then
  # Remove this role and set symlink
  rm -rf "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
  ln -s "${GITHUB_WORKSPACE}" "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
fi

# Provision server
TEST=1 ansible-playbook "${ACTION_PATH}/ansible/server-provision.yml" -i "${INVENTORY_PATH}" -vv
