---
description: Lists all available settings for module configuration files.
---

# Module configuration file

The module configuration file `stackhead-module.yml` is a YAML file that is used for module-specific configuration such as the module type and which Terraform module is required.

{% code title="stackhead-module.yml" %}
```yaml
# this is an example of all available configuration options
---
type: webserver # or container
terraform:
  provider:
    name: myprovider # binary created will be called terraform-provider-myprovider
    url: https://github.com/getstackhead/terraform-myprovider/releases/download/v1.0.0/terraform-provider-myprovider
```
{% endcode %}

## type

The type of your module \(i.e. webserver or container\). While your are currently not required to set this attribute, we recommend you to do so.

{% code title="stackhead-module.yml" %}
```yaml
---
type: webserver
```
{% endcode %}

## terraform

Using the `provider` setting the required Terraform provider can be specified. If set, they will be installed during server setup.

`name` is the actual name of the Terraform provider. Going by Terraform conventions a binary file is called `terraform-provider-[name]` with `[name]` being the actual name required in this setting.

`url` points to the path where the binary file is being downloaded from.

{% code title="stackhead-module.yml" %}
```yaml
---
terraform:
  provider:
    name: caddy # binary created will be called terraform-provider-caddy
    url: https://github.com/getstackhead/terraform-caddy/releases/download/v1.0.0/terraform-provider-caddy
```
{% endcode %}

