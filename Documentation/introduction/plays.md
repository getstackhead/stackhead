---
title: Playbooks
---

The playbooks below are available to be run via `ansible-playbook`:

```shell script
ansible-playbook ansible/playbooks/[file] -i path/to/inventory.yml
```

| File | Description |
| ---- | ----------- |
| server-provision.yml | Perform [server setup](./workflow.md) on all servers in inventory file |
| server-check.yml | Outputs versions of installed software  |
| application-deploy.yml | Perform [project deployment](./workflow.md) for all servers in inventory file |
| application-destroy.yml | Remove all containers and Nginx configurations of a project. Pass in the project name with `--extra-vars "project_name=PROJECTNAME"` |

