# Testing modules

## Integration testing

There is a GitHub action for integration testing, which will:

1. deploy to a server \(specified by `ipaddress` input\)
2. setup a multi-container application and made it available at two domains \(specified by `domain` and `domain2`\)
3. test that content is served correctly
4. tear down the deployed application

In order to use the action you'll have to set up a webserver and make sure GitHub can connect via SSH onto it. Set up two TLDs \(or subdomains\) and point their DNS A record to that IP address.

```yaml
- uses: getstackhead/stackhead@master
  with:
    ipaddress: 'your ip address'
    domain: 'yourdomain.com'
    domain2: 'yourdomain2.com'
    webserver: 'getstackhead.stackhead_webserver_nginx' # webserver to use (make sure to install it)
    container: 'getstackhead.stackhead_container_docker' # container manager to use (make sure to install it)
    plugins: [] # you may define plugins to use for your tests if needed
    rolename: '<vendor>.stackhead_<type>_<name>'
```

