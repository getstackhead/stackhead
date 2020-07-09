#!/bin/bash
# IP address in environment "IP"
# Domain in environment "DOMAIN"
INVENTORY_PATH=ansible/__tests__/inventory.yml

# This test deploys a native project onto the target server
# IMPORTANT: This must run after test_provision.sh!

sed -e "s/\${ipaddress}/${IP}/" -e "s/\${application}/native/" ansible/__tests__/inventory.dist.yml > ansible/__tests__/inventory.yml
sed -e "s/\${domain}/${DOMAIN}/" ansible/__tests__/projects/native.dist.yml > ansible/__tests__/projects/native.yml
TEST=1 ansible-playbook ansible/application-deploy.yml -i $INVENTORY_PATH -vv
content=$(wget --no-check-certificate --http-user=user --http-password=pass --https-only -q -O - https://"${DOMAIN}")
if [[ $content != *"This website was provisioned by StackHead"* ]]; then
  echo "HTTP content check on main domain in native project failed" 1>&2
  exit 1
fi
# test that subdomain is available
content=$(wget --no-check-certificate --https-only -q -O - https://sub."${DOMAIN}":
if [[ $content != *"This website was provisioned by StackHead"* ]]; then
  echo "HTTP content check on subdomain in native project failed" 1>&2
  exit 1
fi
TEST=1 ansible-playbook ansible/application-destroy.yml -i $INVENTORY_PATH --extra-vars "project_name=native" -vv