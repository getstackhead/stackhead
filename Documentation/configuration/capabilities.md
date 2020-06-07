---
title: Capabilities
---

Capabilities are used to define what a target system provides and what a project requires.
Each capability definition requires at least specifying a `version`.

**Server capabilities** are defined in Ansible inventory inside `stackhead.capabilities` of a host definition.
They state what software is additionally installed on the target system during server setup.

**Native-typed projects** can define capabilities to ensure that the target system meets the requirements of the application.
The target server has to fulfill those in order for the application to be deployed.
Deployments of applications where the target system does not fulfill the capabilities (i.e. when the package is missing or has a different version than requested)
will fail.

This page lists all available capabilities and how to configure them.

## PHP

### Inventory definition

```yaml title="Inventory definition"
---
all:
  hosts:
    example:
      stackhead:
        capabilities:
          php:
            version: 7.3
```

### Project definition

```yaml title="Project definition"
---
deployment:
  type: native
  settings:
    capabilities:
      php:  # requires target system with PHP 7.3
        version: 7.3
```