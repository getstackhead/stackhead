#!/bin/bash
# IP address in environment "INPUT_IPADDRESS"

if [[ $INPUT_CLI != '' ]]; then
	CI=1 ${INPUT_CLI_BIN_PATH} setup "${INPUT_IPADDRESS}" -c "/tmp/.stackhead-cli.yml" -v
else
	INVENTORY_PATH="${GITHUB_ACTION_PATH}/inventory.yml"
	sed -e "s/\${ipaddress}/${INPUT_IPADDRESS}/" -e "s/\${dns}/${INPUT_DNS}/" -e "s/\${dns_cloudflare_api_token}/${INPUT_DNS_CLOUDFLARE_APITOKEN}/" -e "s/\${webserver}/${INPUT_WEBSERVER}/" -e "s/\${container}/${INPUT_CONTAINER}/" -e "s/\${plugins}/${INPUT_PLUGINS}/" "${GITHUB_ACTION_PATH}/inventory.dist.yml" >"${INVENTORY_PATH}"

	# Provision server
	CI=1 ansible-playbook "${STACKHEAD_CLONE_LOCATION}/ansible/playbooks/server-provision.yml" -i "${INVENTORY_PATH}" -vv
fi
