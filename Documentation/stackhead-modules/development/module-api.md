---
description: 'Variables, Processes and filters available to use for StackHead plugins.'
---

# Module API

## Variables

The following variables are set by StackHead and can be used by the role:

| Variable | Description | Scope | Data type |
| :--- | :--- | :--- | :--- |
| `stackhead__roles` | Path to local StackHead roles directory | global | string |
| `stackhead__templates` | Path to local StackHead template directory | global | string |
| `stackhead__acme_folder` | Path to remote ACME challenge directory | global | string |
| `stackhead__snakeoil_privkey` | Path to the fake SSL certificate's privkey file | global | string |
| `stackhead__snakeoil_fullchain` | Path to the fake SSL certificate's fullchain file | global | string |
| `stackhead__project_folder` | Path to project's folder | deployment | string |
| `stackhead__tf_project_folder` | Path to project's Terraform folder | deployment | string |
| `stackhead__certificates_project_folder` | Path to project's SSL certificate folder | deployment | string |
| `project_name` | Name of the project that is being deployed | deployment | string |
| `app_config` | Contents of the project definition file | deployment | object |

## Processes

### Terraform execution

If you need to apply your written Terraform configuration within your task, please use the following tasks:

```yaml
- import_tasks: "{{ stackhead__roles }}/stackhead_module_api/tasks/terraform.yml"
```

### Generate SSL certificates

If you want to generate a SSL certificate for the project, run the following task:

```yaml
- include_tasks: "{{ stackhead__roles }}/stackhead_module_api/tasks/ssl-certificate.yml"
```

This will prepare a Terraform configuration file for generating SSL certificates.

You'll find the files inside the project's certificate directory:

* Private key: `{{ stackhead__project_certificates_folder }}/privkey.pem`
* Full chain: `{{ stackhead__project_certificates_folder }}/fullchain.pem`

{% hint style="info" %}
The certificate files will be created after your role was executed. If you need to access the certificate files within your role, please execute Terraform as described above.
{% endhint %}

## Filters

{% hint style="info" %}
All filters provided have to be addressed with the prefix `getstackhead.stackhead`.
{% endhint %}

{% hint style="danger" %}
These filters are pretty specific for the Docker container provider. With the introduction of more container providers these filters may change.
{% endhint %}

### containerPorts

The `containerPorts` filter is used to map the entries from the exposed service configuration \(`containerapp__expose`\) into a more processable format, including the reference to the external port of the Terraform object. The output will be an array of objects with the following structure:

```yaml
{
  # position in output list
  'index': integer,
  # name of exposed service
  'service': string,
  # internal container port that is being exposed
  'internal_port': internal_port,
  # reference to external Docker container port (which is chosen by system)
  'tfstring': "${docker_container.stackhead-" + project_name + "-" + service_name + ".ports[" + str(index) + "].external}"
}
```

Syntax: `containerapp__expose|getstackhead.stackhead.containerPorts(string projectName)`

```yaml
# Given the following expose list:
containerapp__expose:
  - service: pma
    internal_port: 80
    external_port: 81
  - service: nginx
    internal_port: 80
    external_port: 80
  - service: pma
    internal_port: 80
    external_port: 80

# Example usage
{% set all_ports = containerapp__expose|getstackhead.stackhead.containerPorts('myproject') %}

# Output (value of all_ports):

[{
  'index': 0,
  'service': 'nginx',
  'internal_port': 80,
  # reference to external Docker container port (which is chosen by system)
  'tfstring': "${docker_container.stackhead-myproject-nginx.ports[0].external}"
},{
  'index': 1,
  'service': 'pma',
  'internal_port': 80,
  # reference to external Docker container port (which is chosen by system)
  'tfstring': "${docker_container.stackhead-myproject-pma.ports[0].external}"
},{
  'index': 2,
  'service': 'pma',
  'internal_port': 80,
  # reference to external Docker container port (which is chosen by system)
  'tfstring': "${docker_container.stackhead-myproject-pma.ports[1].external}"
}]
```

### TFreplace\(string projectName\)

The `TFreplace` filter is used to replace `$DOCKER_SERVICE_NAME` placeholder variables with the actual reference to the Terraform container object.

Syntax: `"any string"|getstackhead.stackhead.TFreplace(string projectName)`

```yaml
# Example usage
{{ "$DOCKER_SERVICE_NAME['0'] - $DOCKER_SERVICE_NAME['1']"|getstackhead.stackhead.TFreplace('myproject') }}

# Result:

${docker_container.stackhead-myproject-0.name} - ${docker_container.stackhead-myproject-1.name}
```

