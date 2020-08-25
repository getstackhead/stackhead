---
title: Installation
---

# Installation

## Prerequisities

The following software has to already be installed on your system to use StackHead:

* Ansible
* Python 3

## Installing StackHead

StackHead consists of multiple Ansible playbooks which you can install using several package managers.

{% hint style="warning" %}
When cloning from Git make sure to also initialize the submodules by running the command `git submodule update --init` afterwards.
{% endhint %}

### Using Ansible

TBD

### Using Composer

As of right now there is no stable release of StackHead. Please install it directly from the repository as below:

{% code title="composer.json" %}
```javascript
{
    "repositories": [
        { "type": "vcs", "url": "git@github.com:getstackhead/stackhead.git" }
    ],
    "require": {
        "getstackhead/stackhead": "dev-master"
    },
    "scripts":{
        "stackhead-submodules": [
            "cd vendor/getstackhead/stackhead && git submodule sync && git submodule update --init"
        ],
        "post-install-cmd": [
            "@stackhead-submodules"
        ],
        "post-update-cmd": [
            "@stackhead-submodules"
        ]
    }
}
```
{% endcode %}

{% hint style="success" %}
StackHead is installed into `vendor/getstackhead/stackhead` directory.
{% endhint %}

### Using NPM

TBD

