# Playbooks

The playbooks below are available to be run via `ansible-playbook`:

```bash
ansible-playbook $ANSIBLE_COLLECTION_PATH/getstackhead/stackhead/playbooks/[file] -i path/to/inventory.yml
```

{% hint style="info" %}
`$ANSIBLE_COLLECTION_PATH` refers to the location where your Ansible collections are being installed to. Per default this should be `~/.ansible/collections/ansible_collections/`.
{% endhint %}

| File | Description |
| :--- | :--- |
| setup-ansible.yml | Install all additional dependency required to use the other Ansible playbooks |
| server-provision.yml | Perform [server setup](https://github.com/getstackhead/stackhead/tree/24507df46a619397da3abc9fcf3cd4ef9f9fbdbf/Documentation/technical-documentation/workflow.md) on all servers in inventory file |
| server-check.yml | Outputs versions of installed software |
| application-deploy.yml | Perform [project deployment](https://github.com/getstackhead/stackhead/tree/24507df46a619397da3abc9fcf3cd4ef9f9fbdbf/Documentation/technical-documentation/workflow.md) for all servers in inventory file |
| application-destroy.yml | Remove all containers and Nginx configurations of a project. Pass in the project name with `--extra-vars "project_name=PROJECTNAME"` |

