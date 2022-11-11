---
description: This page describes what StackHead modules are.
---

# About modules

StackHead is organized in components which are interchangeable. They can be configured by setting the respective setting in CLI config file.

## Proxy

Proxy modules provide configuration in order to set up a webserver to use as reverse proxy onto containers.

{% hint style="info" %}
Set the webserver you want to use via `modules.proxy` setting in your `.stackhead-cli.yml` file.
{% endhint %}

## Container managers

Container managers are applications like Docker that provide container technologies.

{% hint style="info" %}
Set the container manager you want to use `modules.container` setting in your `.stackhead-cli.yml` file.
{% endhint %}

## Plugins

Plugins are additional applications to be installed on your server.
Such may include reverse proxies, load balancers, databases, etc.

{% hint style="info" %}
Set the plugins you want to use in the `modules.plugins` array setting in your `.stackhead-cli.yml` file.
{% endhint %}

