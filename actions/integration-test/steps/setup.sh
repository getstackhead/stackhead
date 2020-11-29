#!/bin/bash
# IP address in environment "INPUT_IPADDRESS"

if [[ $INPUT_CLI != '' ]]; then
	${INPUT_CLI_BIN_PATH} project setup "${INPUT_IPADDRESS}" -c "/tmp/.stackhead-cli.yml" -v
else
	INVENTORY_PATH="${GITHUB_ACTION_PATH}/inventory.yml"
	sed -e "s/\${ipaddress}/${INPUT_IPADDRESS}/" -e "s/\${webserver}/${INPUT_WEBSERVER}/" -e "s/\${container}/${INPUT_CONTAINER}/" -e "s/\${plugins}/${INPUT_PLUGINS}/" "${GITHUB_ACTION_PATH}/inventory.dist.yml" >"${INVENTORY_PATH}"

	# Provision server
	TEST=1 ansible-playbook "${STACKHEAD_CLONE_LOCATION}/ansible/playbooks/server-provision.yml" -i "${INVENTORY_PATH}" -vv
fi
