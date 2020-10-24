rtmpauth
========
_rtmpauth_ is an authentication & notification system to be used along side the nginx rtmp module.

## Background
_rtmpauth_ project was built to serve the needs of a small Discord community to allow high quality video streaming to a private rtmp server to remain social during the COVID-19 pandemic. When a member starts streaming, a notification is posted in discord as well as when the stream gains or loses a viewer.  
Each member may also have a twitch channel configured in _rtmpauth_ which differs from their discord/publisher user name. When a twitch channel is defined for a member, notifications will be posted when the twitch stream is live/off-line. This functionality can also serve as general twitch notifications for favorite streamers of the discord community.

_rtmpauth_ was started in May 2020 then made public and developed open source throughout October 2020 as part Hacktoberfest. With many open source projects being assaulted with useless _spam_ contributions in attempt to cheat a free Hacktoberfest t-shirt, I wanted to contribute something this year that exceeds the minimum _"4 lines of code/documentation"_ requirement. Hopefully in the future Hacktoberfest will find a way determine quality contributions to open source, but for this year I hope this effort is valid and recognized.

## Features
- Simple authentication for NGiNX RTMP module
- Discord channel notifications via webhook
- Twitch stream notifications
- HTTP REST user management
- Embedded database
- Single binary deployment

## Installation
Installation documentation WIP

TL;DR - compile project to a binary and either setup as a service with systemd or deploy the project in a container. Ensure all environment variables are configured in the next section of this document.

## Configuration
The project is configured with environment variables. All required exported variables are provided with defaults in `init/rtmpauth.env`.

1. Create a local copy of the file
    ```
    mkdir .local
    cp init/rtmpauth.env .local/rtmpauth.env
    ```
2. Update the variables to suit your needs

3. Source the file
    ```
    source .local/rtmpauth.env
    ```

## Managing RTMP Publishers
User management can be performed with some basic REST calls. You can build a custom application around the API or you can simply interact with via your favorite REST client. For the sake of simplicity, the following examples will be demonstrated using the `curl` command.
NOTE: The primary key for all records in database is the publisher _name_ (discord user name)

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
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o rtmpauth main.go
```

## Security considerations
While it is possible to run this service on a different host, it is intended to run on the same host/container pod as nginx and communicate via localhost. Due to this assumption, the _rtmpauth_ service should NOT be publicly accessible or firewall rules should be configured to only allow connection from the nginx host/container.
