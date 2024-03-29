# Development

## Requirements

* Go 1.18

## Build

Build the binary by running the build script at `.build/build.sh`.

## Testing

When developing for StackHead we encourage you to test with an actual remote server.

We recommend a basic Ubuntu server on [Hetzner Cloud](https://hetzner.cloud/?ref=n7H3qhWcZ2QS).
Right now the cheapest option comes in at 4,15€ per month (3,56€ Server + 0,60€ IPv4).
However it is charged per-use. So you'll only paying the time the server is actually running.
So you should be paying only a few cents (or even nothing) when running it for a few hours while testing.

Make sure to setup the server with SSH key access, so you can connect to it from your local PC with root user.
Verify you can connect to it via `ssh root@[IPv4 address]`.

Then, set the A record of an actual domain or subdomain to the IP address.

Setup server:
`./bin/stackhead-cli setup [IPv4 address]`

Deploy project:
`./bin/stackhead-cli project deploy my_file.stackhead.yml [IPv4 address]`
