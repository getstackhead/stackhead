---
title: Capabilities
---

This page lists all available capabilities and how to configure them.

Capabilities can be inside **web-based** project definitions to ensure that the target system meets the requirements of the application.

## PHP

### Inventory definition

```yaml
---
all:
  hosts:
    example:
      stackhead:
        capabilities:
          php:
            version: 7.3
```

### Project definition

```yaml
---
deployment:
  type: web
  settings:
    capabilities:
      php:
        version: 7.3
```