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

* **webserver**: Configuration for reverse proxy webserver
* **container**: Configuration for launching containers

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

## Structure requirements

StackHead modules are included as role during setup and deployment processes.
As StackHead modules combine logics for both setup and configuration, the variable `stackhead_action` is used to tell the role what it should do.
Please make sure your role executes the correct tasks for each action.

| `stackhead_action` value | Expected behaviour                                                                  | Used in step       |
| ------------------------ | ----------------------------------------------------------------------------------- | ------------------ |
| setup                    | The software is installed.                                                          | Server setup       |
| deploy                   | The software is configured for the given project.                                   | Project deployment |
| load-config              | Load stackhead-module.yml vars into variable given by *include_varname* (see below) | ? |

### Recommended structure

We recommend setting up task files named like the `stackhead_action` value (`tasks/[stackhead_action value].yml`)
and include them in the `tasks/main.yml` like so:

```yaml
---
- include_tasks: "{{ role_path }}/tasks/{{ stackhead_action }}.yml"
```

The `load-config` action is needed as we can't determine the role location.
So we need a task that assigns the contents of the _stackhead-module.yml_ to where StackHead wants it.
Going by the recommended structure, set this inside your `tasks/load-config.yml` file:

```yaml
---
- include_vars:
    file: "{{ role_path }}/stackhead-module.yml"
    name: "{{ include_varname }}"
```

## StackHead module API

### Variables

The following variables are set by StackHead and can be used by the role:

| Variable                                 | Description                                         | Scope      | Data type |
| ---------------------------------------- | --------------------------------------------------- | ---------- | --------- |
| `stackhead__roles`                       | Path to local StackHead roles directory             | global     | string    |
| `stackhead__templates`                   | Path to local StackHead template directory          | global     | string    |
| `stackhead__snakeoil_privkey`            | Path to the fake SSL certificate's privkey file     | global     | string    |
| `stackhead__snakeoil_fullchain`          | Path to the fake SSL certificate's fullchain file   | global     | string    |
| `stackhead__project_folder`              | Path to project's folder                            | deployment | string    |
| `stackhead__tf_project_folder`           | Path to project's Terraform folder                  | deployment | string    |
| `stackhead__certificates_project_folder` | Path to project's SSL certificate folder            | deployment | string    |
| `project_name`                           | Name of the project that is being deployed          | deployment | string    |
| `app_config`                             | Contents of the project definition file             | deployment | object    |

### Terraform execution

If you need to apply your written Terraform configuration within your task, please use the following tasks:

```yaml
- import_tasks: "{{ stackhead__roles }}/stackhead_module_api/tasks/terraform.yml"
```

### Generate SSL certificates

If you want to generate a SSL certificate for the project, run the following task:

```yaml
- include_tasks: "{{ stackhead__roles }}/stackhead_module_api/tasks/ssl-certificate.yml"
```

This will prepare a Terraform configuration file for generating SSL certificates.

You'll find the files inside the project's certificate directory:

* Private key: `{{ stackhead__project_certificates_folder }}/privkey.pem`
* Full chain: `{{ stackhead__project_certificates_folder }}/fullchain.pem`

:::note
The certificate files will be created after your role was executed.
If you need to access the certificate files within your role, please execute Terraform as described above.
:::
