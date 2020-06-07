---
title: Getting started
---

This guide will explain how to **setup** a server for usage with StackHead and **deploy** a basic Docker-based application.

You will require:

* StackHead installed (see [Installation Guide](./installation.md))
* a top level domain
* a webserver with SSH root access

:::note   
In this guide `<LOCAL_STACKHEAD_PATH>` will refer to the location where you installed StackHead.
This path usually depends on how you installed StackHead.   
:::

## File structure

_my-inventory.yml_ will be the Ansible inventory file we're using.
Ultimately, our project definitions are stored in _stackhead_ directory.

That means we have the following file structure:

* my-inventory.yml
* stackhead
  * example_app.yml
  
:::note   
Per default StackHead looks for project definition files in the _stackhead_ directory which is in the same directory as the inventory file.   
:::

## Creating a project defintion

Create a new project definitions file at `./stackhead/example_app.yml` and the following content:

```yaml
---
domain: example.com
deployment:
  type: container
  settings:
    expose:
      service: app
      port: 80
    services:
      - name: app
        image: nginxdemos/hello:latest
```

This defines that a new Docker container shall be created with the `nginxdemos/hello:latest` image.
When `example.com` is opened, the request shall be redirected (proxied) to the container's port 80, which is where the
image's internal web server runs.

Make sure to replace _example.com_ with your own domain.

## Provisioning the server

Let's create an inventory file and provision our first server.
We recommend provisioning only newly created servers to minify the side effects of already installed packages.

### Inventory file

Create a new file e.g. `my-inventory.yml` and set the following content:

```yaml
all:
  vars:
    ansible_user: root
    ansible_connection: ssh
  hosts:
    myhost:
      ansible_host: 123.456.789.10
      stackhead:
        applications:
          - example_app
```

Make sure to replace the `123.456.789.10` with your own server IP.

Operations on the target server are done with `ssh` and the `root` user.
Make sure you have a SSH certificate and added the private key on your server.

Looking at the `stackhead.applications` section we specified the name of our application "example_app" we created earlier.
This means our application will be deployed to that IP.

### Server provisioning

Now we should be ready to go to provision our server.

Make sure you have all Ansible dependencies installed using:

```shell script
ansible-galaxy install -r <LOCAL_STACKHEAD_PATH>/ansible/requirements/requirements.yml
```

Then run the following command to provision your server:

```shell script
ansible-playbook <LOCAL_STACKHEAD_PATH>/ansible/server-provision.yml -i my-inventory.yml
```

### Deploying the project

Before deploying the project, check your domain's DNS settings.
Make sure the A record points to the server IP, as this is required for SSL certificate generation.

Then deploy the project with:

```shell script
ansible-playbook <LOCAL_STACKHEAD_PATH>/ansible/application-deploy.yml -i my-inventory.yml
```

After deployment, open the domain in your web browser.
It should display content and have a valid SSL certificate.

### Destroying the project

Now let's remove all configurations we created during deployment.

```shell script
ansible-playbook <LOCAL_STACKHEAD_PATH>/ansible/application-destroy.yml -i my-inventory.yml --extra-vars "project_name=example_app"
```