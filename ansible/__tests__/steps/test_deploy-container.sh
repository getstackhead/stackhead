#!/bin/bash
# IP address in environment "IP"
# Domain in environment "DOMAIN"
INVENTORY_PATH=ansible/__tests__/inventory.yml

# This test deploys a container project onto the target server
# IMPORTANT: This must run after test_provision.sh!

sed -e "s/\${ipaddress}/${IP}/" -e "s/\${webserver}/${WEBSERVER}/" -e "s/\${application}/container/" ansible/__tests__/inventory.dist.yml > ansible/__tests__/inventory.yml
sed -e "s/\${domain}/${DOMAIN}/" ansible/__tests__/projects/container.dist.yml > ansible/__tests__/projects/container.yml
TEST=1 ansible-playbook ansible/application-deploy.yml -i $INVENTORY_PATH -vv
URL = "https://${DOMAIN}"
content=$(wget --no-check-certificate --https-only -q -O - "${URL}")
if [[ $content != *"Hello world!"* ]]; then
  echo "HTTP content check on container project ($DOMAIN) failed: $content" 1>&2
  exit 1
fi
# test that phpmyadmin is available
URL = "https://${DOMAIN}:81"
content=$(wget --no-check-certificate --https-only -q -O - "${URL}")
if [[ $content != *"phpMyAdmin"* ]]; then
  echo "HTTP content check on phpmyadmin in container project failed: $content" 1>&2
  exit 1
fi
# test that phpmyadmin is available on subdomain
URL = "https://sub.${DOMAIN}"
content=$(wget --no-check-certificate --http-user=user --http-password=pass --https-only -q -O - "${URL}")
if [[ $content != *"phpMyAdmin"* ]]; then
  echo "HTTP content check on subdomain phpmyadmin in container project failed: $content" 1>&2
  exit 1
fi
TEST=1 ansible-playbook ansible/application-destroy.yml -i $INVENTORY_PATH --extra-vars "project_name=container" -vv