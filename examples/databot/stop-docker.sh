#!/bin/bash

cd "$(dirname "$0")"

[[ -w /var/run/docker.sock ]] && sudo="" || sudo="sudo"

exec sh -c "${sudo} docker stop databot && ${sudo} docker rm databot"
