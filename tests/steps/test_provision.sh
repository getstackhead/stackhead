#!/bin/bash
# IP address in environment "INPUT_IPADDRESS"
# Domain in environment "DOMAIN"
INVENTORY_PATH=./tests/inventory.yml

sed -e "s/\${ipaddress}/${INPUT_IPADDRESS}/" -e "s/\${webserver}/${INPUT_WEBSERVER}/" ./tests/inventory.dist.yml > $INVENTORY_PATH

# Install dependencies
ansible-galaxy install -r ./ansible/requirements/requirements.yml --force-with-deps

if [[ $INPUT_ROLENAME != '' ]]; then
  # Remove this role and set symlink
  rm -rf "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
  ln -s "$(pwd)" "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
fi

cat "${INVENTORY_PATH}"

# Provision server
TEST=1 ansible-playbook ./ansible/server-provision.yml -i "${INVENTORY_PATH}" -vv
