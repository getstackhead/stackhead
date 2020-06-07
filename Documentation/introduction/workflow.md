---
title: Workflow
---

StackHead utilizes Ansible and Terraform to set up your projects.

![stackhead-process]

The figure above illustrates the general StackHead workflow.
StackHead provides custom Ansible playbooks for installing required software on a remote server (setup) and configuring your projects (deployment).

The highlighted terms are explained in further detail below.

[stackhead-process]: /img/docs/StackHead-process.png "StackHead process"

## Project definition

A web project usually includes an application and additional components (runtime environment, databases, etc).
Also it is usually served on a Top Level Domain by some kind of webserver.

This information is stored inside a [**project definition** file](../configuration/project-definition.md).

There you can either set up container-based applications (recommended) or also run you application natively.

Based on the project definition file, StackHead will take care setting up the required configurations.

## Ansible inventory

Ansible is used to connect to the remote server.
The information on the remote servers are stored in a inventory file.
For more information on Ansible inventories in general, please have a look at the [official Ansible documentation](https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html).

For each host you can define which projects should run there. Also, using [capabilities](../configuration/capabilities.md) you can install additional software environments during server setup.

## Server setup

During server setup all software and utilities that are required to set up your projects with StackHead are installed.
Such include Terraform, Docker and Nginx. You'll have to run the server setup before you can deploy projects onto it.

The server setup is executed by running the respective Ansible playbook.

## Project deployment

Setting up a project is called deployment and is done by running the respective Ansible playbook.
This will create all project-related resources such as Nginx configuration, SSL certificates and start up containers.

:::note  
Only servers that have been set up using the StackHead **server setup** can be deployed onto!  
:::

## Resource management

StackHead uses Terraform for creating resources of any kind.
Such include configuration files, SSL certificates and Docker containers.

While understanding how we use Terraform is not required for using StackHead,
you can find out more in the [technical documentation](../technical-documentation/workflow.md) if you are interested.