# Installation

## Prerequisities

The following software has to already be installed on your system to use StackHead:

* Ansible \(&gt;= 2.10\)
* Python 3

## Installing StackHead

Get the latest StackHead binary from our GitHub repository.

If you wish to configure the StackHead modules to use for webservers or container,
create a `.stackhead-cli.yml` file either in the working directory or in the user home directory:

```yaml
---
modules:
  webserver: nginx
  container: docker
```

{% hint style="info" %}
Modules are resolved automatically to StackHead namespace, e.g. the webserver value `nginx` is treated as `getstackhead.stackhead_webserver_nginx`.
If you're not using an official StackHead module, please make sure to add the vendor name (e.g. `acme.myserver`).
{% endhint %}

Then run the following command to install the StackHead Ansible collection and dependencies.

```bash
stackhead-cli init
```
