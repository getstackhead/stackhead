---
description: >-
  Project definitions contain everything that is needed to set up your web
  application.
---

# Project definition

Project definitions are stored at `./stackhead/[projectname].stackhead.yml`. Each file consists of a **domain** and an **application configuration**.

Applications run in one or multiple Docker containers. While the settings are based on Docker Compose version 2.4, some may require a different syntax. Please have a closer look at the list below.

{% hint style="info" %}
The file has to end in **.stackhead.yml** or **.stackhead.yaml**.
Otherwise the deployment will fail!
{% endhint %}

{% hint style="warning" %}
Not all options from Docker Compose are supported right now.
{% endhint %}

## Annotated project definition

The file below is a fully annotated project definition file.

```yaml
---
domains:
  - domain: example.com # domain name (without protocol and port)
    expose: # expose one or multiple ports
      - internal_port: 80 # port inside the container
        external_port: 80 # port where service can be reached by browsers (i.e. example.com:80)
        service: app # name of service that should be exposed
    security:
      authentication:
        - type: basic # basic authentication: users will have to authenticate with username ("user") and password ("pass")
          username: user
          password: pass
container:
  services:
    - name: app # service name
      image: nginxdemos/hello:latest # Docker image name
      volumes:
        - type: global # mount the "config" folder inside container at "/var/config". As both services use "global" they will mount the same source and share data.
          src: config
          dest: /var/config
        - type: local # mount the "data" folder inside container at "/var/data". This folder is not shared with other services ("local").
          src: data
          dest: /var/data
        - type: custom # mount the "/docker/data/test" folder inside container at "/var/test".
          src: /docker/data/test
          dest: /var/test
      hooks:
        execute_after_setup: ./setup.sh # this file is executed inside container after the container is created
        execute_before_destroy: ./destroy.sh # this file is executed inside container before the container is destroyed
    - name: db # service name
      image: mariadb:10.5 # Docker image name
      environment: # environment variables for Docker container
        MYSQL_ROOT_PASSWORD: example
      volumes_from: app:ro # mount all volumes from app service as readonly
```

## Settings

### domains.\*.expose

The web server will proxy all web traffic to the service and port specified in `expose` setting.

In the example above, the web server will proxy web requests to the "app" container's port 80.

#### service

Name of the Container service to receive the web request.

#### internal\_port

Port of the given container service to receive the web request.

#### external\_port

Port that Nginx listens to.

{% hint style="danger" %}
Setting _external\_port_ to 443 is not allowed, as HTTPS forwarding is automatically enabled for exposes with `external_port=80`.
{% endhint %}

{% hint style="warning" %}
Make sure to define the different _external\_port_ within one project definition, so that each port is only used once!
{% endhint %}

### domains.\*.security

These options can be used to add further security to your projects.

#### authentication

The authentication setting with `type: basic` will require the user to log in with a name and password. You may specify how many users you like.

Removing the `authentication` section will remove the file containing the usernames and passwords for your project.

```yaml
domains:
  - domain: mydomain.com
    security:
      authentication:
        - type: basic
          username: user1
          password: pass1
        - type: basic
          username: user2
          password: pass2
```

{% hint style="info" %}
Right now, removing a single entry from the list and redeploying the project will NOT remove the user settings from the authentication file.
{% endhint %}

### container.services

The following configuration options are available inside a service definition:

#### name

Internal name of your service. Used as service name in generated docker-compose file.

#### image \(required\)

See [docker-compose documentation on image](https://docs.docker.com/compose/compose-file/compose-file-v2/#image)

#### hooks

```yaml
services:
  - name: nginx
    # ...
    hooks:
      execute_after_setup: ./setup.sh
      execute_before_destroy: ./destroy.sh
```

{% hint style="info" %}
Only the specified file is copied onto the server (and container) during deployment.
Note that you should not use or call any files in the file that are not already on the container!
{% endhint %}

#### volumes

StackHead saves mounted data in the project directory at project or service level. You can also define a custom location on the server.

| Configuration | Description | Allowed values |
| :--- | :--- | :--- |
| `type` \(required\) | Determines the data storage location | "global", "local" or "custom" |
|  | **global**: Data storage location is located at `/stackhead/projects/[project_name]/container_data/global/` |  |
|  | **local**: Data storage location is located at `/stackhead/projects/[project_name]/container_data/services/[service_name]/` |  |
|  | **custom**: No data storage location. You have to set it yourself using the _src_ setting below \(absolute path!\). |  |
| `src`  \(required for type=custom\) | Relative path inside data storage location that should be mounted.  Note: When type=custom this is has to be an absolute path! | any string |
| `dest` | Absolute path inside the Docker container where the mount should be applied | any string |
| `mode` | Defines if the volume should be read-write \(rw\) or readonly \(ro\) | "rw" \(default\) or "ro" |

Below you can see a comparison of the project definition \(left\) and the equivalent docker-compose definition:

{% tabs %}
{% tab title="StackHead" %}
{% code title="example\_project.yml" %}
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
{% endcode %}
{% endtab %}

{% tab title="Docker Compose" %}
{% code title="docker-compose.yml" %}
```yaml
services:
  nginx:
    # ...
    volumes:
      - /stackhead/projects/example_project/container_data/global/assets:/var/www/public/assets:rw
      - /stackhead/projects/example_project/container_data/services/nginx/log:/var/www/public/log:rw
      - /etc/secrets.txt:/var/www/secrets.txt:ro
```
{% endcode %}
{% endtab %}
{% endtabs %}

#### volumes\_from

See [docker-compose documentation on volumes\_from](https://docs.docker.com/compose/compose-file/compose-file-v2/#volumes_from).

#### environment

See [docker-compose documentation on environment](https://docs.docker.com/compose/compose-file/compose-file-v2/#environment).

#### user

See [docker-compose documentation on user](https://docs.docker.com/compose/compose-file/compose-file-v2/#user).

