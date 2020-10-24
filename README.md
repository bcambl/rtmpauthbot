rtmpauth
========
`rtmpauth` is an authentication & notification system to be used allong side the nginx rtmp module. While it is possible to run this service on a different host, it is intended to run on the same host as nginx and communicate via localhost.

## Features
- Simple authentication for NGiNX RTMP module
- Discord channel notifications via webhook
- Twitch stream notifications
- HTTP REST user management
- Embedded database
- Single binary deployment

## Installation
WIP

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

#### Adding a publisher:


## Publisher Twitch Stream Motifications

## Build From Source

## Notes
