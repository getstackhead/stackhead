# Configuration

StackHead will look for a file named `.stackhead-cli.yml` in the working directory or in the home directory of the user executing the command.

This file is used to configure which StackHead modules to use.

### Full annotated configuration

```yaml
---
modules:
  webserver: nginx
  container: docker
```

{% hint style="info" %}
Modules are resolved automatically to StackHead namespace, e.g. the web server value `nginx` is treated as `getstackhead.stackhead_webserver_nginx`. If you're not using an official StackHead module, please make sure to add the vendor name \(e.g. `acme.stackhead_webserver_nginx`\).
{% endhint %}

