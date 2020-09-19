---
description: Lists all available settings for module configuration files.
---

# Module configuration file

The module configuration file `stackhead-module.yml` is a YAML file that is used for module-specific configuration such as the module type and which Terraform module is required.

{% code title="stackhead-module.yml" %}
```yaml
# this is an example of all available configuration options
---
constraints:
  - stackhead>=1.0.0
type: webserver # or container
terraform:
  provider:
    name: myprovider # binary created will be called terraform-provider-myprovider
    url: https://github.com/getstackhead/terraform-myprovider/releases/download/v1.0.0/terraform-provider-myprovider
```
{% endcode %}

## constraints

Constraints are used to ensure a StackHead module is compatible with other components.
You can define which versions of StackHead the module is compatible with.

The contraint format follows [semantic versioning](https://semver.org).

```yaml
---
constraints:
  - stackhead>=1.0.0 # All releases greater or equal than version 1.0.0
  - stackhead ^1.0.0 # All minor and patch releases (>=1.0.0 <2.0.0)
  - stackhead ~1.0.0 # All patch releases (>=1.0.0 <1.1.0)
  - stackhead >=1.3.0,<2.1 # All releases from 1.3.0 up to 2.1.*
```

{% hint style="info" %}
We recommend allowing any release after the version that is known to work (e.g. stackhead>=1.0.0).
Once a version is known to be incompatible, the constraint can be adjusted.
This is to make sure modules do not have to be updated once a new (breaking) StackHead is published.
{% endhint %}

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

The provider is installed from Terraform registry, setting `vendor`, `name` and `version`.

`name` is the actual name of the Terraform provider. `vendor` is the owner's name on Terraform registry.
Looking at the provider _getstackhead/caddy_, _getstackhead_ is the vendor and _caddy_ is the name.

{% code title="stackhead-module.yml" %}
```yaml
---
terraform:
  provider: # source=getstackhead/caddy
    vendor: getstackhead
    name: caddy
    version: 1.0.1
```
{% endcode %}

