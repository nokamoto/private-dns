# private-dns

[![Build Status](https://travis-ci.org/nokamoto/private-dns.svg?branch=master)](https://travis-ci.org/nokamoto/private-dns)
[![CircleCI](https://circleci.com/gh/nokamoto/private-dns/tree/master.svg?style=svg)](https://circleci.com/gh/nokamoto/private-dns/tree/master)

## Install dnscli

```
$ make
$ dnscli
```

## Run Private DNS

```
docker run --name private-dns -p 53:53/tcp -p 53:53/udp -p 9999:9999/tcp --cap-add=NET_ADMIN nokamotohub/private-dns
```
