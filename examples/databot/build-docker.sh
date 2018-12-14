#!/bin/bash

cd "$(dirname "$0")"

[[ -w /var/run/docker.sock ]] && sudo="" || sudo="sudo"

exec ${sudo} docker build -t databot .
