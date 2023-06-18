# Docker Container module

![Maintained](https://img.shields.io/badge/status-maintained-green)

## About this module

This is a module to manage Containers via Docker.

Credentials for private registries are encrypted and securely storaged via [Pass](https://www.passwordstore.org/).

**Important note**: The credentials are grouped per registry and saved per username, so only one password is assigned to the pair (registry, username).
For example, this means that you can only store one passwort for your GitLab.com account named "ciuser".
It is recommended to use one username per project. When using GitLab's project-based Access Tokens, the username can be freely set.

## Resources

Setting up the server with this plugin will install the following:

* Docker CE
* Docker Compose Plugin
* containerd.io
* Pass (for registry credential storage)
  * **Note:** Right now defining multiple credentials for the same registry results in authentication issues. That is why the registry credentials are only available during deployment and not persisted on the target server.
* golang-docker-credential-helpers

## Configuration

None.

## Troubleshooting

### I can not connect to my private Docker Registry

Make sure to configure the private Docker registry credentials in your project definition file.

```yaml
container:
  registries:
    - username: gitlab
      password: glpat-12345abcdeFGHIJ12345
      url: https://registry.gitlab.com
```

**GitLab Note:** Access Tokens require the `read_registry` scope and at least the Reporter role.
