---
title: StackHead modules
---

## Specification

A StackHead module is an Ansible role that contains all steps and templates to:

1. Install and setup the software during server setup
2. Create configurations for this software during project deployment via Terraform

While the [regular Ansible role directory layout](https://docs.ansible.com/ansible/latest/user_guide/playbooks_best_practices.html#directory-layout) apply,
the role has also to adhere to this specification.

### Module types

There are the following module types:

* webserver
** Configuration for reverse proxy webserver
* container
** Configuration for launching containers

### Role name

A role name has to adhere to this schema: `stackhead_[type]_[name]`.

### Required files

Each StackHead module is to required to have a module configuration file (see below) in its root directory.

## Module configuration file

The module configuration file `stackhead-module.yml` is a YAML file that can have the following settings.

### terraform

Using the `provider` setting the required Terraform provider can be specified.
If set, they will be installed during server setup.

`name` is the actual name of the Terraform provider.
Going by Terraform conventions a binary file is called `terraform-provider-[name]` with `[name]` being the actual name required in this setting.

`url` points to the path where the binary file is being downloaded from.

```yaml
---
terraform:
  provider:
    name: caddy # binary created will be called terraform-provider-caddy
    url: https://github.com/getstackhead/terraform-caddy/releases/download/v1.0.0/terraform-provider-caddy
```

## How StackHead uses modules

StackHead modules are included as role during setup and deployment process.
The variable `stackhead_action` will correspond to the respective action that is performed (`setup` or `deploy`).
Please make sure your role executes the correct tasks for each action.

We recommend creating two separate task files:

* `tasks/setup.yml`: Instructions that will install the software.
* `tasks/deploy.yml`: Instructions that will configure the software for a project.
* `tasks/load-config.yml`: Load stackhead-module.yml vars into variable given by *include_varname*
    ```yaml
    ---
    - include_vars:
        file: "{{ role_path }}/stackhead-module.yml"
        name: "{{ include_varname }}"
    ```

Then in your `tasks/main.yml`, simply include them based on the action:

```yaml
---
- include_tasks: "{{ role_path }}/tasks/{{ stackhead_action }}.yml"
```
