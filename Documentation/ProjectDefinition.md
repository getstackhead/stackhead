# Project definitions

Project definitions are stored at `./stackhead/[projectname].yml` (per default).
However that can be overwritten by setting the `stackhead__remote_config_folder` in inventory file.
Each file consists of a **domain** and an **application configuration**.

There are two application types: docker and web. Only one application type is allowed.

## Application Type: docker (recommended)

Docker-based applications are applications that run in one or multiple Docker containers.
The definition is pretty similar to docker-compose. In fact StackHead uses docker-compose to spin up the containers.

**Note:** We're using **version 2.4** of Docker Compose, however not all container options are supported right now.

The example below consists of two services (app and db).

```yaml
---
domain: example.com
deployment:
  type: docker
  settings:
    expose:
      port: 80
      service: app
    services:
      - name: app
        image: nginxdemos/hello:latest
      - name: db
        image: mariadb:10.5
        environment:
          MYSQL_ROOT_PASSWORD: example
```
### expose

The Nginx webserver will proxy all web traffic to the service and port specified in `expose` setting.

In the example above, Nginx will proxy web requests to the "app" container's port 80.

### services

The following configuration options are available inside a service definition:

#### name

Internal name of your service. Used as service name in generated docker-compose file.

#### image (required)

See [docker-compose documentation on image](https://docs.docker.com/compose/compose-file/compose-file-v2/#image)

#### volumes

StackHead saves mounted data in the project directory at project or service level. You can also define a custom location on the server.

| Configuration | Description | Allowed values |
| ------------- | ----------- | -------------- |
| type<br>(required) | Determines the data storage location | "global", "local" or "custom" |
|                 | **global**: Data storage location is located at `/stackhead/projects/[project_name]/docker_data/global/` | |
|                 | **local**: Data storage location is located at `/stackhead/projects/[project_name]/docker_data/services/[service_name]/` | |
|                 |  **custom**: No data storage location. You have to set it yourself using the _src_ setting below (absolute path!). | |
| src <br>(required for type=custom)            | Relative path inside data storage location that should be mounted.<br><br>Note: When type=custom this is has to be an absolute path! | any string |
| dest            | Absolute path inside the Docker container where the mount should be applied | any string |
| mode            | Defines if the volume should be read-write (rw) or readonly (ro) | "rw" (default) or "ro"|


Below you can see a comparison of the project definition (left) and the equivalent docker-compose definition:

<table>
<tr>
<th>Project definition (example_project.yml)</th>
<th>Docker-Compose equivalent</th>
</tr>
<tr>
<td>

```yaml
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

</td>
<td>

```yaml
services:
  nginx:
    # ...
    volumes:
      - /stackhead/projects/example_project/docker_data/global/assets:/var/www/public/assets:rw
      - /stackhead/projects/example_project/docker_data/services/nginx/log:/var/www/public/log:rw
      - /etc/secrets.txt:/var/www/secrets.txt:ro
```

</td>
</tr>
</table>



#### volumes_from

See [docker-compose documentation on volumes_from](https://docs.docker.com/compose/compose-file/compose-file-v2/#volumes_from).

#### environment

See [docker-compose documentation on environment](https://docs.docker.com/compose/compose-file/compose-file-v2/#environment).

#### user

See [docker-compose documentation on user](https://docs.docker.com/compose/compose-file/compose-file-v2/#user).

#### links

See [docker-compose documentation on links](https://docs.docker.com/compose/compose-file/#links).

## Application Type: web

Web-based applications are basic applications that live on the target machine and are served by the Nginx webserver.

```yaml
---
domain: example.com
deployment:
  type: web
```

### Serve location

Per default files are served from the `htdocs` folder inside the `/var/www/[projectname]/` directory.
If you want to serve files from a different folder inside that directory, use `deployment.settings.public_path` as below.

```yaml
---
domain: example.com
deployment:
  type: web
  settings:
    public_path: public
```

### Capabilities

If your application requires other software or runtime environments,
define capabilities using `deployment.settings.capabilities` to make sure the application
is only deployed to targets that meet the requirements.

```yaml
---
domain: example.com
deployment:
  type: web
  settings:
    capabilities:
      php:
        version: 7.3
```

Please also have a look at the [complete list of available capabilities](../Features/Capabilities.md).
