# private-dns

[![Build Status](https://travis-ci.org/nokamoto/private-dns.svg?branch=master)](https://travis-ci.org/nokamoto/private-dns)

## Install dnscli

```
$ make
$ dnscli
```

## Run Private DNS

```
docker run --name private-dns -p 53:53/tcp -p 53:53/udp -p 9999:9999/tcp --cap-add=NET_ADMIN nokamotohub/private-dns
```
