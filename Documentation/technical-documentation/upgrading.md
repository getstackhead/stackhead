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
5. Remove system terraform-providers.tf file: `rm -rf /stackhead/terraform/system/terraform-providers.tf`
6. Symlink the system terraform-providers.tf file: `ln -s /stackhead/terraform/projects/terraform-providers.tf /stackhead/terraform/system/terraform-providers.tf`
