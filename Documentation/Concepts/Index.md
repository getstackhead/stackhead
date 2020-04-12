# Concepts

This page describes the concepts used by Mackerel.

## Projects

Projects consist of your application and a host domain where your application can be reached.  
The information about a project is stored in the **project definition** file.

While we recommend running your application in a Docker container,
you can also set it up to run directly on the target system.
Please have a look at the [project definition documentation](../Configuration/Project.md) for more information.

## Provision & Deploy

When using Mackerel you'll come upon two words:
Server **provisioning** and project **deployment**.

### Server provisioning

Any system you want to use Mackerel on has to be **provisioned** first.
This will install all required software packages and services onto the target system.

Find out how to provision your server in the [Getting started guide](../GettingStarted/Index.md).

### Project deployment

Getting your own application on the provisioned server is called **deployment**.
This will setup the required web server configurations, SSL certificates and application for your project.

Find out how to deploy to your provisioned server in the [Getting started guide](../GettingStarted/Index.md).

## Capabilities

Capabilities are used to define what a target system provides and what a project requires.
Each capability definition requires specifying a `version`.

Please also have a look at the [complete list of available capabilities](../Features/Capabilities.md).

### Server

Server capabilities are defined in Ansible inventory inside `mackerel.capabilities` of a host definition.
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
      mackerel:
        applications:
          - example_com
        capabilities: # installs PHP 7.3
          php:
            version: 7.3
```

### Project

**Web-based** projects can define the same capabilities which the target server has to fulfill in order for the application to be deployed.

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