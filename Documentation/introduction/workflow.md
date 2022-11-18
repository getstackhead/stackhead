# Workflow

Outdated process image:
![StackHead process](../.gitbook/assets/stackhead-process.png)

The figure above illustrates the general StackHead workflow.
StackHead provides tooling for installing required software on a remote server \(setup\) and configuring your projects \(deployment\).

The highlighted terms are explained in further detail below.

## Project definition

A web project usually includes an application and additional components \(runtime environment, databases, etc\). Also it is usually served on a Top Level Domain by some kind of webserver.

This information is stored inside a [**project definition** file](project-definition.md).

Based on the project definition file, StackHead will take care setting up the required configuration. It will start up the specified Docker containers and set up the required web server configuration.

## Server setup

During server setup all software and utilities that are required to set up your projects with StackHead are installed.
Such may include container management software \(e.g. Docker\) and web server software \(e.g. Nginx\). You'll have to run the server setup before you can deploy projects onto it.

## Project deployment

Setting up a project is called deployment.
A deployment will create all project-related resources such as reverse proxy configuration, SSL certificates and it will start up containers.

{% hint style="info" %}
Only servers that have been set up using the StackHead **server setup** can be deployed onto!
In the future it will be possible to deploy onto servers not provisioned by StackHead (see https://github.com/getstackhead/stackhead/issues/169).
{% endhint %}
