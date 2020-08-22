#!/bin/bash
# IP address in environment "IP"
# Domain in environment "DOMAIN"
INVENTORY_PATH=ansible/__tests__/inventory.yml

# map webserver name to role name
if [ "$WEBSERVER" == 'nginx' ]; then WEBSERVER='getstackhead.stackhead_webserver_nginx'; fi
if [ "$WEBSERVER" == 'caddy' ]; then WEBSERVER='stackhead_webserver_caddy'; fi

sed -e "s/\${ipaddress}/${IP}/" -e "s/\${webserver}/${WEBSERVER}/" ansible/__tests__/inventory.dist.yml > $INVENTORY_PATH

# Install dependencies
STACKHEAD_COLLECTION_PATH=~/.ansible/collections/ansible_collections/getstackhead/stackhead
pip install -r "$STACKHEAD_COLLECTION_PATH"/requirements/pip.txt

# Install requirements
ansible-galaxy install -r "$STACKHEAD_COLLECTION_PATH"/requirements/requirements.yml
ansible-galaxy install getstackhead.stackhead_webserver_nginx

# Provision server
TEST=1 ansible-playbook  "$STACKHEAD_COLLECTION_PATH"/playbooks/server-provision.yml -i "${INVENTORY_PATH}" -vv
