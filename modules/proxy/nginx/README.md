# Nginx proxy module

![Unmaintained](https://img.shields.io/badge/status-unmaintained-red)

## About this module

This module sets up a Nginx webserver on the target server, which will
then be used to accept web requests and forward/proxy it to the container
of the actual application behind the domain.

There also is a logic for certificate generation via Certbot.

## Installation

Running `stackhead setup` with this module enabled will install:

* Nginx
* Certbot

## Configuration

The following settings are available in CLI config:

```yaml
modules_config:
  nginx:
    certificates_email: "certificates@saitho.me" # mail address used for certificate renewals
    config:
      User: "www-data"
      ConfPath  : "/etc/nginx/conf.d"
      VhostPath: "/etc/nginx/sites-enabled"
      ErrorLog  : "/var/log/nginx/error.log"
      AccessLog: "/var/log/nginx/access.log"
      PidFile   : "/run/nginx.pid"
      WorkerProcesses: "auto"
      ExtraConfOptions: # additional settings for nginx.conf
        foo: bar
      ExtraConfHttpOptions: # additional settings for nginx.conf inside http block
        foo: bar
      WorkerConnections: 1024
      MultiAccept: off
      MimeFilePath: /etc/nginx/mime.types
      ServerNamesHashBucketSize: 64
      ClientMaxBodySize: "64m"
      Sendfile: "on"
      TcpNopush: "on"
      TcpNodelay: "on"
      ServerTokens: "on"
      ProxyCachePath: ""
      KeepaliveTimeout: 65
      KeepaliveRequests: 100
      TypesHashMaxSize: 2048
```
