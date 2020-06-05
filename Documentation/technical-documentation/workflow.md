---
title: Terraform
---

StackHead generates Terraform configuration files based on the project definition.

The Terraform configuration files are placed in the `terraform` folder inside the project directory.

The Terraform configuration files from all project directories are symlinked into a general Terraform folder (as seen below).
With that Terraform is able to manage resources for all projects simultaneously.

![tf-apply]
![ansible-terraform]


### Execution

TODO: Image with process, describe how TF files are being generated and executed

[tf-apply]: /img/docs/terraform-files-structure.png "Applying Terraform changes"
[ansible-terraform]: /img/docs/ansible-terraform-interaction.png "StackHead Workflow: Ansible and Terraform"