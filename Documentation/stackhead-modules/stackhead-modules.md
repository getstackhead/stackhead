---
description: This page describes what StackHead modules are.
---

# About modules

StackHead is organized in components which are interchangable. They can be configured by setting the respective variable in Ansible inventory definition.

## Webservers

Webserver modules provide configuration in order to set up a webserver software to use as reverse proxy onto containers.

{% hint style="info" %}
Set the webserver you want to use via `stackhead__webserver` in your Ansible inventory file.
{% endhint %}

## Container managers

Container managers are applications like Docker that provide container technologies.

{% hint style="info" %}
Set the container manager you want to use via `stackhead__container` in your Ansible inventory file.
{% endhint %}

