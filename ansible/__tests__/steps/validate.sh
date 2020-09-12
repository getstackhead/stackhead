#!/bin/bash
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
ping -c 5 "${INPUT_DOMAIN}"
ping -c 5 "${INPUT_DOMAIN2}"

ssl_check "${INPUT_DOMAIN}"
ssl_check "${INPUT_DOMAIN2}"
http_check "${INPUT_DOMAIN}" "Hello world!"
http_check "${INPUT_DOMAIN}:81" "phpMyAdmin"
http_check "${INPUT_DOMAIN2}" "phpMyAdmin" "user" "pass"
