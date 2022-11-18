---
description: This page describes upgrading to StackHead v2.
---

# Upgrading to StackHead v2

We recommend setting up the projects from scratch. However it should be possible to upgrade projects deployed with StackHead v1 to v2 via manual migration (see below).

## Breaking Changes

### Removed Terraform

Terraform is nice and great, but it is a software that requires updates. The main software and plugins as well.
There is currently no concept in StackHead for upgrading Terraform components, therefore it was dropped altogether.

The resources are now created directly on the file system:

* Nginx: TBD
* Caddy: TBD
* Docker: TBD

### Modules

#### Destroy

All modules must implement a `Destroy` method, called when a project is teared down.

### StackHead CLI configuration

```diff
---
modules:
-  webserver: nginx
+  proxy: nginx
   container: docker
   dns:
     - cloudflare
-  plugins:
-    - watchtower # load getstackhead.stackhead_plugin_watchtower plugin

-certificates:
-  register_email: "my-certificates-mail@mydomain.com" # Email address used for creating SSL certificates. Will receive notice when they expire.
-config:
-  setup:
-    github.com/getstackhead/plugin-proxy-nginx: # config settings for Nginx module
-      server_names_hash_bucket_size: 128
-      extra_conf_options:
-        foo: bar
-      extra_conf_http_options:
-        foo2: bar2
+ modules_config:
+   nginx:
+     certificates_email: "my-certificates-mail@mydomain.com"
+     config:
+       server_names_hash_bucket_size: 128
+       extra_conf_options:
+         foo: bar
+       extra_conf_http_options:
+         foo2: bar2
```
