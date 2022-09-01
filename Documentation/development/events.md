---
description: Events within StackHead
---

This page describes the events used in StackHead.

Plugins may listen onto these events to do additional things.

Example:
```go
import "github.com/gookit/event"

// executed after docker module was set up (i.e. Docker is ready to use on the system)
event.On("setup.modules.post-install-module.container.docker", event.ListenerFunc(func(e event.Event) error {

})
```

{% hint style="info" %}
As of right now, plugins are only able to listen to events that occur after they have been installed.
{% endhint %}

# Setup

The following events are used during the setup step:

| Event Name                                       | Description                                                                                                                                                                  | Event params       |
|--------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------|
| `setup.folders.pre-install`                      | Triggered before StackHead folders are set up                                                                                                                                | -                  |
| `setup.folders.post-install`                     | Triggered after StackHead folders are set up                                                                                                                                 | -                  |
| `setup.users.pre-install`                        | Triggered before StackHead users are set up                                                                                                                                  | -                  |
| `setup.users.post-install`                       | Triggered after StackHead users are set up                                                                                                                                   | -                  |
| `setup.ssh.pre-install`                          | Triggered before StackHead SSH keys are set up                                                                                                                               | -                  |
| `setup.ssh.post-install`                         | Triggered after StackHead SSH keys are set up                                                                                                                                | -                  |
| `setup.modules.pre-install`                      | Triggered before StackHead modules are set up                                                                                                                                | -                  |
| `setup.modules.post-install`                     | Triggered after StackHead modules are set up                                                                                                                                 | -                  |
| `setup.modules.pre-install-module.[type].[name]` | Triggered before a _specific_ StackHead modules install step is executed.<br/>`[type]` = module type (e.g. proxy, container)<br/>`[name]` = module name (e.g. caddy, docker) | `{"module": Module}` |
| `setup.modules.post-install-module.[type].[name]` | Triggered after a _specific_ StackHead modules install step is executed.<br/>`[type]` = module type (e.g. proxy, container)<br/>`[name]` = module name (e.g. caddy, docker)  | `{"module": Module}` |

{% hint style="info" %}
The order of which modules are installed is: container, proxy
{% endhint %}
