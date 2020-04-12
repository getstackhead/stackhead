# deployment

## Dependencies

Make sure to install Ansible dependencies via:
```
ansible-galaxy install -r requirements/requirements.yml
```

## Workflows

For local deployment, run workflow with `-i inventories/local.yaml -K`.
For remote deployment, run workflow with `-i inventories/remote.yaml`.

| File | Description |
| ---- | ----------- |
| server-provision.yml | Setup required programs on server  |
| server-check.yml | Outputs versions of installed software  |
| application-deploy.yml | Setup application containers and Nginx config |

Please checkout the [Documentation/Index.md](documentation) and
the [Getting started guide](Documentation/GettingStarted/Index.md) for more information on usage.