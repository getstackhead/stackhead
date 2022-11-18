# Configuration

StackHead will look for a file named `.stackhead-cli.yml` in the working directory or in the home directory of the user executing the command.

This file is used to configure which StackHead modules to use.

You may define additional module configurations within the `modules_config` key.
See example below for setting the setting _server_names_hash_bucket_size_ for the Nginx proxy module.

{% hint style="warning" %}
The `modules`, `modules_config` and `terraform` settings may be overwritten via server configuration.
The server configuration is located in `/stackhead/config.yml` (if it exists).
{% endhint %}

## Full annotated configuration

```yaml
---
modules:
  proxy: nginx
  container: docker
  plugins:
    - portainer
  dns:
    - cloudflare
modules_config:
  nginx: # config settings for Nginx module
    certificates_email: "my-certificates-mail@mydomain.com" # Email address used for creating SSL certificates. Will receive notice when they expire.
    config:
      server_names_hash_bucket_size: 128
terraform:
  update_interval: "*-*-* 4:00:00" # perform Terraform update everyday at 4am, see Unix timer "OnCalendar" setting
```

{% hint style="info" %}
Please look at the individual README files of the modules for all available configuration settings.
{% endhint %}

