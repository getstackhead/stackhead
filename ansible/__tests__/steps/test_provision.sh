#!/bin/bash
# IP address in environment "IP"
# Domain in environment "DOMAIN"
INVENTORY_PATH=ansible/__tests__/inventory.yml

# map webserver name to role name
if [ "$WEBSERVER" == 'nginx' ]; then WEBSERVER='getstackhead.stackhead_webserver_nginx'; fi
if [ "$WEBSERVER" == 'caddy' ]; then WEBSERVER='stackhead_webserver_caddy'; fi

sed -e "s/\${ipaddress}/${IP}/" -e "s/\${webserver}/${WEBSERVER}/" ansible/__tests__/inventory.dist.yml > $INVENTORY_PATH

# Install dependencies
pip install -r ansible/requirements/pip.txt
ansible-galaxy install -r ansible/requirements/requirements.yml

# Provision server
TEST=1 ansible-playbook ansible/server-provision.yml -i "${INVENTORY_PATH}" -vv
