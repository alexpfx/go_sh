#!/bin/bash

docker compose build
#docker rmi $(docker images -f "dangling=true" -q)
