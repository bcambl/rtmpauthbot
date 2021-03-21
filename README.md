rtmpauthbot
=========

`rtmpauthbot` is an authentication & notification system to be used along side the nginx rtmp module.

## Background
`rtmpauthbot` was built to provide a small Discord community with a private rtmp server with a simple authentication mechanism and basic notifications for when members start and stop streaming.

Each member may also have a twitch channel configured which differs from their discord/publisher user name. When a twitch channel is defined for a member, notifications will be posted when the twitch stream is live/off-line. This functionality can also serve as general twitch notifications for favorite streamers of the discord community.

## Features
- Authentication system for NGiNX RTMP module
- Discord channel notifications
- Twitch stream notifications
- HTTP REST user management
- Embedded database
- Single binary deployment

## Configuration
The project is configured with environment variables.

1. Create a local copy of the environment variable file
    ```
    mkdir /etc/rtmpauthbot
    rtmpauthbot -environment > /etc/rtmpauthbot/rtmpauthbot.env
    ```
2. Update the variables to suit your needs

## Install Service
Installation documentation WIP

TL;DR - compile project to a binary and either setup as a service with `systemd` or deploy the project in a container. Ensure all environment variables are configured from the previous section of this document.


A basic systemd unit-file can be generated with the following command
```
rtmpauthbot -unitfile > /etc/systemd/system/rtmpauthbot.service
systemctl daemon-reload
```

## Managing RTMP Publishers
User management can be performed with some basic REST calls. You can either interact with `rtmpauthbot` using your favorite REST client or build a custom application around the API. For the sake of simplicity, the following examples will be demonstrated using the `curl` command.  

### Adding/Updating a publisher
```
curl -X POST -d '{"name": "discord_username", "key": "private_rtmp_stream_key"}' http://127.0.0.1:9090/api/publisher
```
expected response status code: `204`

Optionally, If a user would also like to provide notifications for their public twitch stream:
```
curl -X POST -d '{"name": "discord_username", "key": "private_rtmp_stream_key", "twitch_stream": "twitch_username"}' http://127.0.0.1:9090/api/publisher
```
expected response status code: `204`

### Retrieve all publishers
```
curl http://127.0.0.1:9090/api/publisher
```

expected response status code: `200`
```
[
  {
    "name": "discord_username",
    "key": "abcdefghijklmnopqrstuvwxyz0123456789",
    "twitch_stream": "twitch_username"
  }
]
```

### Retrieve a single publisher
```
curl http://127.0.0.1:9090/api/publisher?name=discord_username
```

expected response status code: `200`
```
{
    "name": "discord_username",
    "key": "abcdefghijklmnopqrstuvwxyz0123456789",
    "twitch_stream": "twitch_username"
}
```

### Deleting a publisher
```
curl -X DELETE -d '{"name": "discord_username"}' http://127.0.0.1:9090/api/publisher
```

expected response status code: `204`

## Build From Source
If you would rather compile the project from source, please install the latest version of the Go programming language  [here](https://golang.org/dl/).
```
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o rtmpauthbot main.go
```

## Security considerations
While it is possible to run this service on a different host, it is intended to run on the same host/container pod as nginx and communicate via localhost. Due to this assumption, the `rtmpauthbot` service should NOT be publicly accessible or firewall rules should be configured to only allow connection from the nginx host/container.
