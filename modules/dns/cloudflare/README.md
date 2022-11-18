# Cloudflare DNS module

![Maintained](https://img.shields.io/badge/status-maintained-green)

## About this module

This module allows automatically configuring the DNS settings for a domain.

## Resources

None.

## Configuration

### API Token

Make sure to provide the API token for Cloudflare.
You can generate an API token [in your Cloudflare profile](https://dash.cloudflare.com/profile/api-tokens).
Make sure to grant `write` permissions to DNS on Zone level.

```yaml
modules_config:
  cloudflare:
    # scoped Cloudflare API token
    api_token: MY-API-TOKEN
    # switching safemode off will – during project destroy – remove all DNS entries on the domain
    disable_safemode: false
```

### Domain setting

You'll also have to define the DNS provider to be used for each domain you want to set up in project definition:

```yaml
domains:
  - domain: mydomain.com
    dns:
      provider: cloudflare
```

