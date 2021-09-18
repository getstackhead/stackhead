#!/bin/bash

# Write config file
echo "---
modules:
	dns:
		- ${INPUT_DNS}
	webserver: ${INPUT_WEBSERVER}
	container: ${INPUT_CONTAINER}
	plugins: ${INPUT_PLUGINS}
config:
	setup:
		getstackhead.stackhead_webserver_nginx:
			server_names_hash_bucket_size: 128
	deployment:
		getstackhead.stackhead_dns_cloudflare:
			cloudflare_api_token: ${INPUT_DNS_CLOUDFLARE_APITOKEN}
	destroy:
		getstackhead.stackhead_dns_cloudflare:
			cloudflare_api_token: ${INPUT_DNS_CLOUDFLARE_APITOKEN}
" >"/tmp/.stackhead-cli.yml"

if [[ $INPUT_SELFTEST != '' ]]; then
	${INPUT_CLI_BIN_PATH} init --version="${INPUT_VERSION}" -v -c "/tmp/.stackhead-cli.yml" -v
else
	${INPUT_CLI_BIN_PATH} init --version=next -v -c "/tmp/.stackhead-cli.yml" -v
fi
