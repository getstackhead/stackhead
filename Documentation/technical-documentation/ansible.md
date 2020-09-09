# Ansible

StackHead uses Ansible to install the required software and connect to the servers.
The collection is available on Ansible Galaxy: [getstackhead/stackhead](https://galaxy.ansible.com/getstackhead/stackhead).

It can also be installed from the official Git repository.

```bash
ansible-galaxy collection install getstackhead.stackhead
ansible-galaxy collection install git+https://github.com/getstackhead/stackhead.git
```

## Ansible inventory

Ansible is used to connect to the remote server. The information on the remote servers are stored in a inventory file. For more information on Ansible inventories in general, please have a look at the [official Ansible documentation](https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html).

For each host you can define which projects should run there.

### Variables

Aside from specifying the server IP with the `ansible_host` property, the applications that should be deployed have to be set in `stackhead.applications` setting as below.

The application name is the name of the project definition file without YAML file extension.

Note that you also must specify which StackHead modules to use for webserver and container (`stackhead__webserver` and `stackhead__container` setting).

```yaml
all:
  vars:
    ansible_user: root
    ansible_connection: ssh
  hosts:
    myhost:
      stackhead__webserver: getstackhead.stackhead_webserver_nginx # use NGINX as webserver
      stackhead__container: getstackhead.stackhead_container_docker # use Docker for containers
      ansible_host: 123.456.789.10
      stackhead:
        applications:
          - example_app
```

## Playbooks

The playbooks below are available to be run via `ansible-playbook`:

```bash
ansible-playbook $ANSIBLE_COLLECTION_PATH/getstackhead/stackhead/playbooks/[file] -i path/to/inventory.yml
```

{% hint style="info" %}
`$ANSIBLE_COLLECTION_PATH` refers to the location where your Ansible collections are being installed to.
Per default this should be `~/.ansible/collections/ansible_collections/`.
{% endhint %}

| File | Description |
| :--- | :--- |
| server-provision.yml | Perform [server setup](workflow.md) on all servers in inventory file |
| server-check.yml | Outputs versions of installed software |
| application-deploy.yml | Perform [project deployment](workflow.md) for all servers in inventory file |
| application-destroy.yml | Remove all containers and Nginx configurations of a project. Pass in the project name with `--extra-vars "project_name=PROJECTNAME"` |

