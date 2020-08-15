#!/bin/bash
# Domain in environment "DOMAIN"
# IMPORTANT: This must run after test_deploy.sh!

function http_check() {
  echo "Checking HTTP content on ${1}"
  if [[ "${3}" != "" && "${4}" != "" ]]; then
    CONTENT=$(curl --insecure -v -u "${3}:${4}" "${1}")
  else
    CONTENT=$(curl --insecure -v "${1}")
  fi
  if [[ $CONTENT != *"${2}"* ]]; then
    echo "HTTP content check failed: ${CONTENT}" 1>&2
    exit 1
  else
    echo "HTTP content check succeeded"
  fi
}

http_check "https://${DOMAIN}" "Hello world!"
http_check "https://${DOMAIN}:81" "phpMyAdmin"
http_check "https://sub.${DOMAIN}" "phpMyAdmin" "user" "pass"
