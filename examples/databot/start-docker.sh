#!/bin/bash

cd "$(dirname "$0")"

[[ -w /var/run/docker.sock ]] && sudo="" || sudo="sudo"

exec ${sudo} docker run -d --name databot --restart always --env-file config/keybase.env -v $(pwd)/config:/databot/config -v databot_downloads:/databot/downloads databot:latest

