#!/bin/bash
ansible-galaxy collection build -f "${GITHUB_ACTION_PATH}/ansible"
ansible-galaxy collection install "$(find getstackhead-stackhead-*)" -f

# Install dependencies
if [[ $INPUT_SELFTEST != '' ]]; then
  ansible-galaxy install -r "${GITHUB_WORKSPACE}/ansible/requirements/requirements.yml" --force-with-deps
else
  ansible-galaxy install -r "${GITHUB_ACTION_PATH}/ansible/requirements/requirements.yml" --force-with-deps
fi

if [[ $INPUT_ROLENAME != '' ]]; then
  # Remove this role and set symlink
  rm -rf "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
  ln -s "${GITHUB_WORKSPACE}" "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
fi
