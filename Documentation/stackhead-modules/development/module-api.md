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

