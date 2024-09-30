#!/usr/bin/env bash

export CATTLE_BOOTSTRAP_PASSWORD=$(openssl rand -hex 12)

docker rm -f rancher
docker run -d --restart=unless-stopped -p 8989:80 -p 8943:443 --privileged \
  --name rancher -e CATTLE_BOOTSTRAP_PASSWORD \
  rancher/rancher:latest

echo "Rancher started"
echo "Bootstrap Password: $CATTLE_BOOTSTRAP_PASSWORD"
echo "Wait a minute and then access Rancher at https://localhost:8943"
