#!/bin/bash
# IP address in environment "IP"
# Domain in environment "DOMAIN"
# Ansible inventory path in environment "INVENTORY_PATH"

# This test deploys a container project onto the target server
# IMPORTANT: This must run after test_provision.sh!

# map webserver name to role name
if [ $WEBSERVER == 'nginx' ]; then WEBSERVER='getstackhead.stackhead_module_webserver_nginx'; fi
if [ $WEBSERVER == 'caddy' ]; then WEBSERVER='stackhead_webserver_caddy'; fi

sed -e "s/\${ipaddress}/${IP}/" -e "s/\${webserver}/${WEBSERVER}/" -e "s/\${application}/container/" ansible/__tests__/inventory.dist.yml > ansible/__tests__/inventory.yml
sed -e "s/\${domain}/${DOMAIN}/" ansible/__tests__/projects/container.dist.yml > ansible/__tests__/projects/container.yml
TEST=1 ansible-playbook ansible/application-deploy.yml -i "${INVENTORY_PATH}" -vv
