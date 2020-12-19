#!/bin/bash

if [[ $INPUT_CLI != '' ]]; then
	# Write config file
	echo "---
modules:
	webserver: ${INPUT_WEBSERVER}
	container: ${INPUT_CONTAINER}
	plugins: ${INPUT_PLUGINS}
config:
	setup:
		getstackhead.stackhead_webserver_nginx:
			server_names_hash_bucket_size: 128
" >"/tmp/.stackhead-cli.yml"

	if [[ $INPUT_SELFTEST != '' ]]; then
		${INPUT_CLI_BIN_PATH} init --version="${INPUT_VERSION}" -v -c "/tmp/.stackhead-cli.yml" -v
	else
		${INPUT_CLI_BIN_PATH} init --version=next -v -c "/tmp/.stackhead-cli.yml" -v
	fi
else
	cp "${STACKHEAD_CLONE_LOCATION}/VERSION" "${STACKHEAD_CLONE_LOCATION}/ansible/VERSION"
	rm -rf "${STACKHEAD_CLONE_LOCATION}/ansible/schema"
	rm -rf "${STACKHEAD_CLONE_LOCATION}/ansible/schema"
	mkdir "${STACKHEAD_CLONE_LOCATION}/ansible/schema"
	(cd "${STACKHEAD_CLONE_LOCATION}/schemas" && find . -name "*.json" -exec cp --parents -R '{}' "../ansible/schemas/" ';')
	ansible-galaxy collection build -f "${STACKHEAD_CLONE_LOCATION}/ansible"
	ansible-galaxy collection install "$(find getstackhead-stackhead-*)" -f

	# Install dependencies
	if [[ $INPUT_SELFTEST != '' ]]; then
		pip install -r "${STACKHEAD_CLONE_LOCATION}/ansible/requirements/pip_requirements.txt"
		ansible-galaxy install -r "${STACKHEAD_CLONE_LOCATION}/ansible/requirements/requirements.yml" --force-with-deps
		ansible-playbook "${STACKHEAD_CLONE_LOCATION}/ansible/playbooks/setup-ansible.yml"
	else
		pip install -r "${STACKHEAD_CLONE_LOCATION}/ansible/requirements/pip_requirements.txt"
		ansible-galaxy install -r "${STACKHEAD_CLONE_LOCATION}/ansible/requirements/requirements.yml" --force-with-deps
		ansible-playbook "${STACKHEAD_CLONE_LOCATION}/ansible/playbooks/setup-ansible.yml"
	fi

	ansible-galaxy install "${INPUT_WEBSERVER}"
	ansible-galaxy install "${INPUT_CONTAINER}"
fi

if [[ $INPUT_ROLENAME != '' ]]; then
	# Remove this role and set symlink
	rm -rf "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
	ln -s "${GITHUB_WORKSPACE}" "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
fi
