#!/usr/bin/env bash

set -xe

docker build -t denismakogon/go-plugin-func:latest . -f Dockerfile.test
docker run --rm  denismakogon/go-plugin-func:latest
docker rmi -f denismakogon/go-plugin-func | true
