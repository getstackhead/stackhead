---
title: Container
---

Container-based applications are applications that run in one or multiple Docker containers.
While the settings are based on Docker Compose version 2.4, some may require a different syntax.
Please have a closer look at the list below.

:::note
Not all options from Docker Compose are supported right now.
:::

The example below consists of two services (app and db).

```yaml
---
type: container
domains:
  - domain: example.com
    expose:
      - internal_port: 80
        external_port: 80
        service: app
container:
  services:
    - name: app
      image: nginxdemos/hello:latest
    - name: db
      image: mariadb:10.5
      environment:
        MYSQL_ROOT_PASSWORD: example
```

## domains.*.expose

The Nginx webserver will proxy all web traffic to the service and port specified in `expose` setting.

In the example above, Nginx will proxy web requests to the "app" container's port 80.

### service

Name of the Container service to receive the web request.

### internal_port

Port of the given container service to receive the web request.

### external_port

Port that Nginx listens to.

:::caution  
Setting _external_port_ to 443 is not allowed, as HTTPS forwarding is automatically enabled for exposes with `external_port=80`.  
:::

:::note  
Make sure to define the different _external_port_ within one project definition, so that each port is only used once!  
:::

## container.services

The following configuration options are available inside a service definition:

### name

Internal name of your service. Used as service name in generated docker-compose file.

### image (required)

See [docker-compose documentation on image](https://docs.docker.com/compose/compose-file/compose-file-v2/#image)

### volumes

StackHead saves mounted data in the project directory at project or service level. You can also define a custom location on the server.

| Configuration | Description | Allowed values |
| ------------- | ----------- | -------------- |
| type<br/>(required) | Determines the data storage location | "global", "local" or "custom" |
|                 | **global**: Data storage location is located at `/stackhead/projects/[project_name]/container_data/global/` | |
|                 | **local**: Data storage location is located at `/stackhead/projects/[project_name]/container_data/services/[service_name]/` | |
|                 |  **custom**: No data storage location. You have to set it yourself using the _src_ setting below (absolute path!). | |
| src <br/>(required for type=custom)            | Relative path inside data storage location that should be mounted.<br/><br/>Note: When type=custom this is has to be an absolute path! | any string |
| dest            | Absolute path inside the Docker container where the mount should be applied | any string |
| mode            | Defines if the volume should be read-write (rw) or readonly (ro) | "rw" (default) or "ro"|


Below you can see a comparison of the project definition (left) and the equivalent docker-compose definition:

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<Tabs
  defaultValue="stackhead"
  values={[
    { label: 'StackHead', value: 'stackhead', },
    { label: 'Docker Compose', value: 'dockercompose', },
  ]
}>
<TabItem value="stackhead">

```yaml title="example_project.yml"
services:
  - name: nginx
    # ...
    volumes:
      - type: global
        src: assets
        dest: /var/www/public/assets
      - type: local
        src: log
        dest: /var/www/public/log
      - type: custom
        src: /etc/secrets.txt
        dest: /var/www/secrets.txt
        mode: ro
```

</TabItem>
<TabItem value="dockercompose">

```yaml title="docker-compose.yml"
services:
  nginx:
    # ...
    volumes:
      - /stackhead/projects/example_project/container_data/global/assets:/var/www/public/assets:rw
      - /stackhead/projects/example_project/container_data/services/nginx/log:/var/www/public/log:rw
      - /etc/secrets.txt:/var/www/secrets.txt:ro
```

</TabItem>
</Tabs>

### volumes_from

See [docker-compose documentation on volumes_from](https://docs.docker.com/compose/compose-file/compose-file-v2/#volumes_from).

### environment

See [docker-compose documentation on environment](https://docs.docker.com/compose/compose-file/compose-file-v2/#environment).

### user

See [docker-compose documentation on user](https://docs.docker.com/compose/compose-file/compose-file-v2/#user).
