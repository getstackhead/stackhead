---
description: >-
  This guide will explain how to setup a server deploy a basic Docker-based
  application with StackHead CLI.
---

# Getting started

You will require:

* StackHead CLI binary \(see [GitHub releases page](https://github.com/getstackhead/stackhead/releases)\)
* a top level domain
* a web server with SSH root access

If you wish to change the software used for proxy or containers, please [create a CLI configuration file](cli-configuration.md).

## Creating a project definition

Create a new project definitions file at `./stackhead/example_app.yml` and the following content:

```yaml
---
domains:
  - domain: example.com
    expose:
      - service: app
        external_port: 80
        internal_port: 80
container:
  services:
    - name: app
      image: nginxdemos/hello:latest
```

This defines that a new Docker container shall be created with the `nginxdemos/hello:latest` image. When `example.com` is opened, the request shall be redirected \(proxied\) to the container's port 80, which is where the image's internal web server runs.

**Important:** Make sure to replace _example.com_ with your own domain.

## Server setup

Let's setup our first server. We recommend only doing that with newly created servers to minimize the side effects of already installed packages.

Run the following command to provision your server \(replace `123.456.789.10` with your own server's IP address\)

```bash
# IPv4
stackhead-cli setup 123.456.789.10

# IPv6
stackhead-cli setup 1234:4567:90ab:cdef::1
```

{% hint style="info" %}
During server setup, a `stackhead` user will be created on the remote server, as well as a SSH keypair which is used for connecting during deployment.
The keys are stored on the executing machine at `~/.config/getstackhead/stackhead/ssh/remotes/[IP Address]`.
{% endhint %}
{% hint style="warning" %}
If you ran the server setup with a IPv4 address and want to deploy to the IPv6 address (and vice-versa), make sure to symlink the SSH keys as well:
`ln -s ~/.config/getstackhead/stackhead/ssh/remotes/[IPv4 address] ~/.config/getstackhead/stackhead/ssh/remotes/[IPv6 address]`
{% endhint %}

## Deploying the project

Before deploying the project, check your domain's DNS settings. Make sure the A record points to the server IP, as this is required for SSL certificate generation.

Then deploy the project with:

```bash
# IPv4
stackhead-cli deploy ./stackhead/example_app.yml 123.456.789.10

# IPv6
stackhead-cli deploy ./stackhead/example_app.yml 1234:4567:90ab:cdef::1
```

After deployment, open the domain in your web browser. It should display content and have a valid SSL certificate.

## Destroying the project

Now let's remove all configurations we created during deployment. This will remove the web server configuration and Docker containers.

```bash
stackhead-cli destroy ./stackhead/example_app.yml 123.456.789.10

# IPv6
stackhead-cli destroy ./stackhead/example_app.yml 1234:4567:90ab:cdef::1
```

