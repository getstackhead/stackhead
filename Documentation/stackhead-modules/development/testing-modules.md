# Testing modules

### Integration testing

There is a GitHub action for integration testing, which will:

* deploy to a server \(specified by `ipaddress` input\)
* setup a multi-container application and made it available at two domains \(specified by `domain` and `domain2`\)
* test that content is served correctly
* tear down the deployed application

In order to use the action you'll have to set up a webserver and make sure GitHub can connect via SSH onto it. Set up two TLDs \(or subdomains\) and point their DNS A record to that IP address.

```yaml
- uses: getstackhead/stackhead@master
  with:
    ipaddress: 'your ip address'
    domain: 'yourdomain.com'
    domain2: 'yourdomain2.com'
    webserver: 'getstackhead.stackhead_webserver_nginx' # webserver to use (make sure to install it)
    rolename: '<vendor>.stackhead_<type>_<name>'
```

