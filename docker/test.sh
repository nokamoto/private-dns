#!/bin/bash

set -ex

docker run -d --name private-dns -p 53:53/tcp -p 53:53/udp -p 9999:9999/tcp --cap-add=NET_ADMIN nokamotohub/private-dns

while ! docker exec -it private-dns nc -z localhost 9999; do
  sleep 0.1
done

[[ -z $(docker exec -it private-dns dnscli get) ]]

! docker exec -it private-dns dnscli -apikey x get

IP=127.0.0.1

HOST=test.host

docker exec -it private-dns dnscli add $IP $HOST

docker exec -it private-dns dnscli get | grep $HOST | grep $IP

docker exec -it private-dns dnscli remove $IP $HOST

[[ -z $(docker exec -it private-dns dnscli get) ]]

docker stop private-dns

docker rm private-dns
