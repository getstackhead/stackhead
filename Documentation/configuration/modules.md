# StackHead Modules

StackHead is organized in components which are interexchangable. They can be configured by setting the respective variable in Ansible inventory definition.

## Webservers

Webserver modules provide configuration in order to set up a webserver software to use as reverse proxy onto containers. You can set a webserver to use by overriding the setting _stackhead\_\_webserver\_ in your inventory file.

### List of webserver modules

* [Nginx \(getstackhead.stackhead\_webserver\_nginx\)](https://github.com/getstackhead/module-webserver-nginx)
* Caddy \(stackhead\_webserver\_caddy\), built-in

## Container managers

Container managers are applications like Docker that provide container technologies.

### List of container managers

* [Docker \(getstackhead.stackhead\_container\_docker\)](https://github.com/getstackhead/module-container-docker)

