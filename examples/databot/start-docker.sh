#!/bin/bash

cd "$(dirname "$0")"

sudo=""
if ! test -w /var/run/docker.sock; then
        sudo="sudo"
fi

exec ${sudo} docker run -d --name databot --restart always --env-file config/keybase.env -v $(pwd)/config:/databot/config databot:latest
