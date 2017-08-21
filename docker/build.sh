#!/bin/bash

set -ex

cd docker

docker build -t nokamotohub/private-dns .

./test.sh

docker push nokamotohub/private-dns
