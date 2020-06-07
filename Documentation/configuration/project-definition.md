---
title: Project definition
---

Project definitions are stored at `./stackhead/[projectname].yml` (per default).
However that can be overwritten by setting the `stackhead__remote_config_folder` in inventory file.
Each file consists of a **domain** and an **application configuration**.

There are two application types: [container](./project-definition/container.md) and [native](./project-definition/native.md). Only one application type is allowed.