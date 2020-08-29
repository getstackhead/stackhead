#!/bin/bash
ansible-galaxy collection build -f ansible
ansible-galaxy collection install getstackhead-stackhead-*

# Install dependencies
ansible-galaxy install -r "${GITHUB_ACTION_PATH}/ansible/requirements/requirements.yml" --force-with-deps

if [[ $INPUT_ROLENAME != '' ]]; then
  # Remove this role and set symlink
  rm -rf "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
  ln -s "${GITHUB_WORKSPACE}" "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
fi
