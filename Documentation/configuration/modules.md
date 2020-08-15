---
title: StackHead modules
---

StackHead is organized in components which are interexchangable.
They can be configured by setting the respective variable in Ansible inventory definition.

## Settings

| Setting                | Description                           | Allowed Values | Default |
| -----                  | ------------------------------------- | -------------- | ------- |
| `_stackhead__webserver`| Webserver to use for reverse proxying | nginx, caddy   | nginx   |
