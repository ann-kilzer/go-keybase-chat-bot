#!/bin/bash

cd "$(dirname "$0")"

sudo=""
if ! test -w /var/run/docker.sock; then
        sudo="sudo"
fi

exec sh -c "${sudo} docker stop databot && ${sudo} docker rm databot"
