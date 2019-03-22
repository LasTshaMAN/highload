#!/usr/bin/env bash

docker build --file prometheus/Dockerfile -t my-prom .
docker build -t common .
