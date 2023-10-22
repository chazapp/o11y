#!/bin/bash
set -e
docker build -t ghcr.io/chazapp/o11y/wall_api:$1 .
docker push ghcr.io/chazapp/o11y/wall_api:$1