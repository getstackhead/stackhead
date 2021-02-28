# Configuration

StackHead will look for a file named `.stackhead-cli.yml` in the working directory or in the home directory of the user executing the command.

This file is used to configure which StackHead modules to use.

You may define additional module configurations for each step (setup, deployment, destroy) within the config key.
Note that you have to set the fully qualified plugin name (e.g. `getstackhead.stackhead_webserver_nginx`).

See example below for setting the setting _server_names_hash_bucket_size_ for the Nginx webserver module in setup step.

## Full annotated configuration

```yaml
---
modules:
  webserver: nginx
  container: docker
  plugins:
    - watchtower # load getstackhead.stackhead_plugin_watchtower plugin
certificates:
  register_email: "my-certificates-mail@mydomain.com" # Email address used for creating SSL certificates. Will receive notice when they expire.
terraform:
  update_interval: "*-*-* 4:00:00" # perform Terraform update everyday at 4am, see Unix timer "OnCalendar" setting
config:
  setup:
    getstackhead.stackhead_webserver_nginx: # config settings for Nginx module
      server_names_hash_bucket_size: 128
```

{% hint style="info" %}
Modules defined in the `modules` section are resolved automatically to StackHead namespace, e.g. the web server value `nginx` is treated as `getstackhead.stackhead_webserver_nginx`. If you're not using an official StackHead module, please make sure to add the vendor name \(e.g. `acme.stackhead_webserver_nginx`\).
{% endhint %}

