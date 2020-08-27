#!/bin/bash
# This test deploys a container project onto the target server
# IMPORTANT: This must run after test_provision.sh!

INVENTORY_PATH=./tests/inventory.yml
sed -e "s/\${ipaddress}/${INPUT_IPADDRESS}/" -e "s/\${webserver}/${INPUT_WEBSERVER}/" -e "s/\${application}/container/" ./tests/inventory.dist.yml > "${INVENTORY_PATH}"
sed -e "s/\${domain}/${INPUT_DOMAIN}/" ./tests/projects/container.dist.yml > ./tests/projects/container.yml
TEST=1 ansible-playbook ./ansible/application-deploy.yml -i "${INVENTORY_PATH}" -vv
