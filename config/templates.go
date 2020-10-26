package config

import (
	"fmt"
)

const (
	license = `
BSD 2-Clause License

Copyright (c) 2020, Blayne Campbell
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
	list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
	this list of conditions and the following disclaimer in the documentation
	and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
	
`

	envVars = `
# path to database file
DATA_PATH=""

# auth server listen ip
AUTH_SERVER_IP="127.0.0.1"

# auth server listen port
AUTH_SERVER_PORT="9090"

# rtmp server fqdn (used for discord private stream links)
RTMP_SERVER_FQDN="stream.mydomain.com"

# rtmp server port (default: 1935)
RTMP_SERVER_PORT="1935"

# enable/disable discord integrations
DISCORD_ENABLED=false

# discord channel webhook
DISCORD_WEBHOOK="https://discordapp.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz1234567890"

# enable/disable twitch integrations
TWITCH_ENABLED=false

# twitch api client id
TWITCH_CLIENT_ID="abcd1234"

# twitch api client secret
TWITCH_CLIENT_SECRET="abcd1234"

# twitch poll rate in seconds
TWITCH_POLL_RATE="60"

`
	systemdUnit = `
[Unit]
Description=rtmpauthd rtmp authentication server

[Service]
EnvironmentFile=/etc/rtmpauthd/rtmpauthd.env
Type=simple
User=nginx
WorkingDirectory=/var/cache/nginx
ExecStart=/usr/local/bin/rtmpauthd

[Install]
WantedBy=multi-user.target

`
)

// PrintLicense simply prints the LICENSE to stdout
func PrintLicense() {
	fmt.Println(license)
}

// PrintEnv simply prints the env vars to stdout
func PrintEnv() {
	fmt.Println(envVars)
}

// PrintSystemDUnit simply prints the systemd unitfile to stdout
func PrintSystemDUnit() {
	fmt.Println(systemdUnit)
}
