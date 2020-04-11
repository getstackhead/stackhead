# deployment

## Workflows

For local deployment, run workflow with `-i inventories/local.yaml -K`.
For remote deployment, run workflow with `-i inventories/remote.yaml`.

| File | Description |
| ---- | ----------- |
| provision-server.yml | Setup required programs on server  |
| deploy-application.yml | Setup application containers and Nginx config |

## Local deployment

```
ansible-galaxy install -r requirements/requirements.yml
ansible-playbook provision-server.yml -i inventories/local.yaml -K
```

Please checkout the [Documentation/Index.md](documentation) and
the [Getting started guide](Documentation/GettingStarted/Index.md) for more information on usage.