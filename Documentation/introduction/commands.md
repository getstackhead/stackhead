---
description: This page lists all available StackHead CLI commands.
---

# Commands

Global Flags

| Flag                | Description                                                                             |
|---------------------|-----------------------------------------------------------------------------------------|
| -c, --config string | path to CLI config file (default is ./.stackhead-cli.yaml or $HOME/.stackhead-cli.yaml) |
| -v, --verbose       | Show more output                                                                        |

## Server provisioning

Should be executed on newly created servers.

This command installs the required software (e.g. Docker and Caddy), create the required folder structure
and the non-root StackHead user that is used for deploying projects.

This command needs to be executed before projects can be deployed onto the target server.

{% hint style="info" %}
Right now it is not possible for StackHead to deal with SSH fingerprints. Please connect to the server yourself via SSH
and accept the fingerprint hashes. Otherwise StackHead CLI is stuck in the connecting phase and can not proceed.
{% endhint %}

### Syntax

```shell
stackhead-cli setup [IPv4 or IPv6 address]
```

### Example

```bash
# IPv4
stackhead-cli setup 123.45.67.8

# IPv6
stackhead-cli setup 1234:4567:90ab:cdef::1
```

## Project Deployment

This will process a project definition file and deploy the application to the target server.

{% hint style="info" %}
The target server needs to have been provisioned using the setup command.
{% endhint %}

### Syntax

```shell
stackhead-cli project deploy [path to project definition] [ipv4 address] [--autoconfirm] [--no-rollback]
```

### Flags

| Flag          | Description                                          |
|---------------|------------------------------------------------------|
| --autoconfirm | Changes will be made without asking for confirmation |
| --no-rollback | Do not rollback on errors (useful for debugging)     |

### Example

```bash
./bin/stackhead-cli project deploy my_file.stackhead.yml 123.45.67.8
```


## Project Deletion

This will remove a project from the target server. The containers will be stopped and all data is removed.

### Syntax

```shell
stackhead-cli project destroy [path to project definition] [ipv4 address]
```

### Example

```bash
./bin/stackhead-cli project destroy my_file.stackhead.yml 123.45.67.8
```

## Config file validation

There are two commands you can use in order to validate StackHead configuration files.

### Syntax

```
# Validate Project definition file
./bin/stackhead-cli project validate [path to project definition]

# Validate StackHead CLI configuration file
./bin/stackhead-cli cli validate [path to cli definition]
```

### Example

```bash
# Validate Project definition file
./bin/stackhead-cli project validate my_file.stackhead.yml

# Validate StackHead CLI configuration file
./bin/stackhead-cli cli validate ~/.stackhead-cli.yaml
```

