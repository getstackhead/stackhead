#!/bin/bash

if [[ $INPUT_CLI != '' ]]; then
  # Write config file
  echo "---
modules:
  webserver: ${INPUT_WEBSERVER}
  container: ${INPUT_CONTAINER}
" > "/tmp/.stackhead-cli.yml"

  if [[ $INPUT_SELFTEST != '' ]]; then
    ${INPUT_CLI_BIN_PATH} init --version="${INPUT_VERSION}" -c "/tmp/.stackhead-cli.yml"
  else
    ${INPUT_CLI_BIN_PATH} init --version=next -c "/tmp/.stackhead-cli.yml"
  fi
else
  rm -rf "${GITHUB_ACTION_PATH}/ansible/schema"
  cp -R "${GITHUB_ACTION_PATH}/validation/schema" "${GITHUB_ACTION_PATH}/ansible"
  ansible-galaxy collection build -f "${GITHUB_ACTION_PATH}/ansible"
  ansible-galaxy collection install "$(find getstackhead-stackhead-*)" -f

  # Install dependencies
  if [[ $INPUT_SELFTEST != '' ]]; then
    ansible-galaxy install -r "${GITHUB_WORKSPACE}/ansible/requirements/requirements.yml" --force-with-deps
  else
    ansible-galaxy install -r "${GITHUB_ACTION_PATH}/ansible/requirements/requirements.yml" --force-with-deps
  fi
fi

if [[ $INPUT_ROLENAME != '' ]]; then
  # Remove this role and set symlink
  rm -rf "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
  ln -s "${GITHUB_WORKSPACE}" "${HOME}/.ansible/roles/${INPUT_ROLENAME}"
fi
