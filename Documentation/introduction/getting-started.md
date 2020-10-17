---
description: >-
  This guide will explain how to setup a server deploy a basic Docker-based
  application with StackHead CLI.
---

# Getting started

You will require:

* StackHead CLI \(see [Installation Guide](../stackhead-cli/installation-1.md)\)
* a top level domain
* a web server with SSH root access

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
stackhead-cli setup 123.456.789.10
```

## Deploying the project

Before deploying the project, check your domain's DNS settings. Make sure the A record points to the server IP, as this is required for SSL certificate generation.

Then deploy the project with:

```bash
stackhead-cli deploy ./stackhead/example_app.yml 123.456.789.10
```

After deployment, open the domain in your web browser. It should display content and have a valid SSL certificate.

## Destroying the project

Now let's remove all configurations we created during deployment. This will remove the web server configuration and Docker containers.

```bash
stackhead-cli destroy ./stackhead/example_app.yml 123.456.789.10
```

