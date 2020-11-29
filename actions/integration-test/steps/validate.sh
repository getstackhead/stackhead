#!/bin/bash
# IMPORTANT: This must run after test_deploy.sh!

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

function service_check() {
	SERVICE_OUTPUT=$(ssh root@"${INPUT_IPADDRESS}" systemctl status "stackhead-apply-terraform.${1}")
	if [[ $SERVICE_OUTPUT != *"Loaded: loaded"* ]]; then
		echo "${1} check (loaded) failed."
		exit 1
	fi
	if [[ $SERVICE_OUTPUT != *"stackhead-apply-terraform.${1}; enabled"* ]]; then
		echo "${1} check (enabled) failed."
		exit 1
	fi
	if [[ ${1} == "timer" ]]; then
		if [[ $SERVICE_OUTPUT != *"Active: active (waiting)"* ]]; then
			echo "${1} check (active) failed."
			exit 1
		fi
	fi
	echo "All ${1} checks succeeded."
}

openssl version
curl -V
ping -c 5 "${INPUT_DOMAIN}"
ping -c 5 "${INPUT_DOMAIN2}"

http_check "${INPUT_DOMAIN}" "Hello world!"
http_check "${INPUT_DOMAIN}:81" "phpMyAdmin"
http_check "${INPUT_DOMAIN2}" "phpMyAdmin" "user" "pass"

service_check "timer"
service_check "service"
