#!/usr/bin/env bash

#docker swarm leave --force
#docker swarm init
docker stack deploy --compose-file docker-compose.yml common


# check that service was deployed
#docker stack services common

# get logs for container running in swarm cluster
#docker service logs common_prometheus

