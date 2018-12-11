#!/bin/bash

cd "$(dirname "$0")"

docker run -d --name databot --restart always --env-file config/keybase.env -v $(pwd)/config:/databot/config databot:latest
