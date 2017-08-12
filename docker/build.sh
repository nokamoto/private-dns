#!/bin/bash

set -ex

cd docker

docker build -t nokamotohub/private-dns .

docker push nokamotohub/private-dns
