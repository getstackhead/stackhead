#!/bin/bash

sed -i "s/location \/ {/location \/ {\n       autoindex on;/" /etc/nginx/conf.d/default.conf
nginx -s reload -c /etc/nginx/nginx.conf
