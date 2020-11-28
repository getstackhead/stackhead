# Inventory

Ansible is used to connect to the remote server. The information on the remote servers are stored in a inventory file. For more information on Ansible inventories in general, please have a look at the [official Ansible documentation](https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html).

For each host you can define which projects should run there.

### Variables

Aside from specifying the server IP with the `ansible_host` property, the applications that should be deployed have to be set in `stackhead.applications` setting as below.

The application name is the name of the project definition file without YAML file extension.

Note that you also must specify which StackHead modules to use for webserver and container \(`stackhead__webserver` and `stackhead__container` setting\).

```yaml
all:
  vars:
    ansible_user: root
    ansible_connection: ssh
    ansible_python_interpreter: /usr/bin/python3
  hosts:
    myhost:
      stackhead__webserver: getstackhead.stackhead_webserver_nginx # use NGINX as webserver
      stackhead__container: getstackhead.stackhead_container_docker # use Docker for containers
      stackhead__plugins: [] # list of plugins
      stackhead__tf_update_realtime: yes # if "no" Terraform updates will be performed via cron (every 5 minutes)
      ansible_host: 123.456.789.10
      stackhead:
        applications:
          - example_app
```

{% hint style="warning" %}
Your inventory must specific Python 3 as `ansible_python_interpreter`!
{% endhint %}
