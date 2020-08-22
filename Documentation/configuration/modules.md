---
title: StackHead modules
---

StackHead is organized in components which are interexchangable.
They can be configured by setting the respective variable in Ansible inventory definition.

## Webservers

Webserver modules provide configuration in order to set up a webserver software to use as reverse proxy onto containers.
You can set a webserver to use by overriding the setting _stackhead__webserver_ in your inventory file.

### Plugin list

* [Nginx (getstackhead.stackhead_module_webserver_nginx)](https://github.com/getstackhead/module-webserver-nginx)
* Caddy (stackhead_webserver_caddy), built-in
