---
title: "Native"
---

Native applications are basic applications that live on the target machine and are served by the Nginx webserver.

```yaml
---
domain: example.com
deployment:
  type: native
```

## Serve location

Per default files are served from the `htdocs` folder inside the `/var/www/[projectname]/` directory.
If you want to serve files from a different folder inside that directory, use `deployment.settings.public_path` as below.

```yaml
---
domain: example.com
deployment:
  type: native
  settings:
    public_path: public
```

## Capabilities

If your application requires other software or runtime environments,
define capabilities using `deployment.settings.capabilities` to make sure the application
is only deployed to targets that meet the requirements.

```yaml
---
domain: example.com
deployment:
  type: native
  settings:
    capabilities:
      php:
        version: 7.3
```

Please also have a look at the [complete list of available capabilities](../capabilities.md).
