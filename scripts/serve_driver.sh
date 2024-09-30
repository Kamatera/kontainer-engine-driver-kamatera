#!/usr/bin/env bash

echo rancher download url from inside the docker container: http://172.17.0.1:8944/kontainer-engine-driver-kamatera
python3 -m http.server 8944
