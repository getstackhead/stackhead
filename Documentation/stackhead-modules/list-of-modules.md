---
description: A list of all available StackHead modules.
---

# List of modules

See module page for maintainence status.

## Proxy

* [Caddy](https://github.com/getstackhead/stackhead/tree/next/modules/proxy/caddy) (recommended)
* [Nginx](https://github.com/getstackhead/stackhead/tree/next/modules/proxy/nginx)

{% hint style="success" %}
If no webserver is configured, **Caddy** will be used per default.
{% endhint %}

## Container Managers

* [Docker](https://github.com/getstackhead/stackhead/tree/next/modules/container/docker) (recommended)

{% hint style="success" %}
If no container manager is configured, **Docker** will be used per default.
{% endhint %}

## DNS

* [Docker](https://github.com/getstackhead/stackhead/tree/next/modules/dns/cloudflare)

{% hint style="info" %}
In order to configure a DNS service for a domain, the DNS provider name needs to be set in the domain settings of project definition as well!
{% endhint %}

## Plugins

* [Portainer](https://github.com/getstackhead/stackhead/tree/next/modules/plugin/portainer)
