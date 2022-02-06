---
description: This page describes upgrading to StackHead v2.
---

# Upgrading to StackHead v2

We recommend setting up the projects from scratch. However it should be possible to upgrade projects deployed with StackHead v1 to v2 via manual migration (see below).

## Breaking Changes

### Terraform: Per-project apply

When applying Terraform plans all deployed projects are considered which results in a long runtime when many projects are deployed.
Furthermore another project could cause issues e.g. when hitting the Docker registry rate limit.

The Terraform operation is now directly executed in the project's Terraform directory (`/stackhead/projects/[name]/terraform`).
That's also where the state file is stored. The providers are saved to `/stackhead/terraform/.terraform` and shared across all projects via `TF_DATA_DIR` setting.

In order to interact with projects deployed with StackHead v1, they need to be migrated as follows:

1. Move `terraform-providers.tf` file up one directory to `/stackhead/terraform/terraform-providers.tf`
2. Remove all Terraform symlinks from `/stackhead/terraform/projects` instead of those with "provider" in its name.
3. Run `terraform apply` to teardown all applications
4. Finally, remove the `/stackhead/terraform/projects` directory
* In each project Terraform directory: `/stackhead/projects/[name]/terraform`
  1. Symlink the terraform-providers.tf file: `ln -s /stackhead/terraform/projects/terraform-providers.tf /stackhead/projects/[name]/terraform/terraform-providers.tf`
  2. Run Terraform apply with TF_DATA_DIR: `TF_DATA_DIR=/stackhead/terraform/.terraform terraform apply -auto-approve`
  3. _Optional:_ Create a new `env.sh` file with `export TF_DATA_DIR="/stackhead/terraform/.terraform"` as content, if you also intend to run Terraform manually. If you do, run `source env.sh` before that to load the environment.
5. Remove system terraform-providers.tf file: `rm -rf /stackhead/terraform/system/terraform-providers.tf`
6. Symlink the system terraform-providers.tf file: `ln -s /stackhead/terraform/projects/terraform-providers.tf /stackhead/terraform/system/terraform-providers.tf`

### StackHead CLI configuration

```diff
---
modules:
-  webserver: nginx
+  proxy: github.com/getstackhead/plugin-proxy-nginx
-  container: docker
+  container: github.com/getstackhead/plugin-container-docker
-  dns: cloudflare
+  dns: github.com/getstackhead/plugin-dns-cloudflare
-  plugins:
-    - watchtower # load getstackhead.stackhead_plugin_watchtower plugin
+  applications:
+    - github.com/getstackhead/plugin-application-watchtower # load getstackhead.stackhead_plugin_watchtower plugin
certificates:
  register_email: "my-certificates-mail@mydomain.com" # Email address used for creating SSL certificates. Will receive notice when they expire.
terraform:
  update_interval: "*-*-* 4:00:00" # perform Terraform update everyday at 4am, see Unix timer "OnCalendar" setting
config:
  setup:
    github.com/getstackhead/plugin-proxy-nginx: # config settings for Nginx module
      server_names_hash_bucket_size: 128
```