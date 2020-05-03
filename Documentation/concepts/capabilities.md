---
title: Capabilities
---

Capabilities are used to define what a target system provides and what a project requires.
Each capability definition requires specifying a `version`.

Please also have a look at the [complete list of available capabilities](../configuration/capabilities.md).

### Server

Server capabilities are defined in Ansible inventory inside `stackhead.capabilities` of a host definition.
They state what software is additionally installed on the target system during provisioning.

Example inventory:
```yaml
---
all:
  vars:
    ansible_user: root
    ansible_connection: ssh
  hosts:
    hetzner1:
      ansible_host: 116.203.211.171
      stackhead:
        applications:
          - example_com
        capabilities: # installs PHP 7.3
          php:
            version: 7.3
```

### Project

**Native-typed** projects can define the same capabilities which the target server has to fulfill in order for the application to be deployed.

Deployments of applications where the target system does not fulfill the capabilities (i.e. when the package is missing or has a different version than requested)
will fail.

Example project definition:
```yaml
---
domain: example.com
deployment:
  type: web
  settings:
      public_path: public
      capabilities: # requires target system with PHP 7.3
        php:
          version: 7.3
```