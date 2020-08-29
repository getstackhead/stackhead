# Installation

## Prerequisities

The following software has to already be installed on your system to use StackHead:

* Ansible (>= 2.10)
* Python 3

## Installing StackHead

StackHead consists of multiple Ansible playbooks which you can install as Ansible collection:

```shell script
ansible-galaxy collection install git+https://github.com/getstackhead/stackhead.git#/ansible/
```

{% hint style="warning" %}
When cloning from Git make sure to also initialize the submodules by running the command `git submodule update --init` afterwards.
{% endhint %}
