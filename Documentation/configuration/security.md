---
title: Security
---

These options can be used to add further security to your projects.
Add these to your project definition file.

## Authentication

Require users to authenticate before they can access your site.

### HTTP Basic auth

Require user to log in with a name and password. You may specify how many users you like.

Removing the `authentication` section will remove the file containing the usernames and passwords for your project.

```yaml
security:
  authentication:
    - type: basic
      username: user1
      password: pass1
    - type: basic
      username: user2
      password: pass2
```

:::note
Right now, removing a single entry from the list and redeploying the project will NOT remove the user settings from the authentication file.
:::
