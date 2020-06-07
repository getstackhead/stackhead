---
title: Installation
---

## Prerequisities

The following software has to already be installed on your system to use StackHead:

* Ansible

## Installing StackHead

StackHead consists of multiple Ansible playbooks which you can install using several package managers.

:::important   
When cloning from Git make sure to also initialize the submodules by running the command `git submodule update --init` afterwards.   
:::

### Using Ansible

TBD

### Using Composer

As of right now there is no stable release of StackHead.
Please install it directly from the repository as below:

```json title="composer.json"
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
:::tip   
StackHead is installed into `vendor/getstackhead/stackhead`.   
:::


### Using NPM

TBD