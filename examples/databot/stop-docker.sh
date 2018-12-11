#!/bin/bash

cd "$(dirname "$0")"

docker stop databot && docker rm databot