#!/bin/bash
# Domain in environment "DOMAIN"
# IMPORTANT: This must run after test_deploy.sh!

function ssl_check() {
  echo "Checking SSL certificate on ${1}"
  openssl s_client -connect "${1}:443" -servername "${1}"
}

function http_check() {
  echo "Checking HTTP content on ${1}"
  if [[ "${3}" != "" && "${4}" != "" ]]; then
    CONTENT=$(curl --insecure -v -u "${3}:${4}" "https://${1}")
  else
    CONTENT=$(curl --insecure -v "https://${1}")
  fi
  if [[ $CONTENT != *"${2}"* ]]; then
    echo "HTTP content check failed: ${CONTENT}" 1>&2
    exit 1
  else
    echo "HTTP content check succeeded"
  fi
}

openssl version
curl -V
ping -c 5 "${DOMAIN}"
ping -c 5 "sub.${DOMAIN}"

ssl_check "${DOMAIN}"
ssl_check "sub.${DOMAIN}"
http_check "${DOMAIN}" "Hello world!"
http_check "${DOMAIN}:81" "phpMyAdmin"
http_check "sub.${DOMAIN}" "phpMyAdmin" "user" "pass"
