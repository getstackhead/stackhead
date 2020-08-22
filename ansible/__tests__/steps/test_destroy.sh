#!/bin/bash
# IP address in environment "IP"
# Domain in environment "DOMAIN"
# Ansible inventory path in environment "INVENTORY_PATH"

# This test destroys a deployed project on the target server
# IMPORTANT: This must run after test_deploy.sh!

STACKHEAD_COLLECTION_PATH=~/.ansible/collections/ansible_collections/getstackhead/stackhead
TEST=1 ansible-playbook "$STACKHEAD_COLLECTION_PATH"/playbooks/application-destroy.yml -i "${INVENTORY_PATH}" --extra-vars "project_name=container" -vv
