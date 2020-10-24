rtmpauth
========
_rtmpauth_ is an authentication & notification system to be used along side the nginx rtmp module.

## Background
The _rtmpauth_ project was built to serve the needs of a small Discord community to allow high quality video streaming to a private rtmp server. When a member starts streaming, a notification is posted in discord as well as when the stream gains or loses a viewer.  
Each member may also have a twitch channel configured in _rtmpauth_ which differs from their discord/publisher username. When a twitch channel is defined for a member, notifications will be posted when the twitch stream is live/off-line. This functionality can also serve as general twitch notifications for favorite streamers of the discord community.

## Features
- Simple authentication for NGiNX RTMP module
- Discord channel notifications via webhook
- Twitch stream notifications
- HTTP REST user management
- Embedded database
- Single binary deployment

## Installation
WIP

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
The primary key for all records in database is the publisher _name_.

#### Retrieve all publishers

#### Adding a publisher

#### Retrieve a publisher

#### Updating a publisher

#### Deleting a publisher

## Publisher Twitch Stream Notifications

## Build From Source

## Security considerations
While it is possible to run this service on a different host, it is intended to run on the same host/container pod as nginx and communicate via localhost. Due to this assumption, the _rtmpauth_ service should NOT be publicly accessible or firewall rules should be configured to only allow connection from the nginx host/container.

## Notes
