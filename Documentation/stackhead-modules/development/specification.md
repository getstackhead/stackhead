# Specification

A StackHead module is an Ansible role that contains all steps and templates to:

1. Install and setup the software during server setup
2. Create configurations for this software during project deployment via Terraform

While the [regular Ansible role directory layout](https://docs.ansible.com/ansible/latest/user_guide/playbooks_best_practices.html#directory-layout) apply, the role has also to adhere to this specification.

## Module types

There are the following module types:

* **webserver**: Configuration for reverse proxy webserver
* **container**: Configuration for launching containers

## Role name

A role name has to adhere to this schema: `stackhead_[type]_[name]`.

## Required files

Each StackHead module is to required to have a [module configuration file](module-configuration-file.md) in its root directory.

## Structure requirements

StackHead modules are included as role during setup and deployment processes. As StackHead modules combine logics for both setup and configuration, the variable `stackhead_action` is used to tell the role what it should do. Please make sure your role executes the correct tasks for each action.

| `stackhead_action` value | Expected behaviour | Used in step |
| :--- | :--- | :--- |
| setup | The software is installed. | Server setup |
| load-config | Load stackhead-module.yml vars into variable given by _include\_varname_ \(see below\) | Server setup |
| deploy | The software is configured for the given project. | Project deployment |

## Recommended structure

We recommend setting up task files named like the `stackhead_action` value \(`tasks/[stackhead_action value].yml`\) and include them in the `tasks/main.yml` like so:

```yaml
---
- include_tasks: "{{ role_path }}/tasks/{{ stackhead_action }}.yml"
```

The `load-config` action is needed as we can't determine the role location. So we need a task that assigns the contents of the _stackhead-module.yml_ to where StackHead wants it. Going by the recommended structure, set this inside your `tasks/load-config.yml` file:

```yaml
---
- include_vars:
    file: "{{ role_path }}/stackhead-module.yml"
    name: "{{ include_varname }}"
```

