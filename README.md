# private-dns

[![Build Status](https://travis-ci.org/nokamoto/private-dns.svg?branch=master)](https://travis-ci.org/nokamoto/private-dns)
[![CircleCI](https://circleci.com/gh/nokamoto/private-dns/tree/master.svg?style=svg)](https://circleci.com/gh/nokamoto/private-dns/tree/master)

## QuickStart
```
$ docker run --name private-dns \
    -p 53:53/tcp -p 53:53/udp -p 9999:9999/tcp \
    --cap-add=NET_ADMIN \
    nokamotohub/private-dns
    
$ dnscli add 192.168.33.10 private-dns.vagrant

# wait for about 60 seconds

$ dig private-dns.vagrant @localhost +short
192.168.33.10
```

### Install dnscli

```
$ make
$ dnscli
```

### Run Private DNS

```
$ docker run --name private-dns \
    -p 53:53/tcp -p 53:53/udp -p 9999:9999/tcp \
    --cap-add=NET_ADMIN \
    nokamotohub/private-dns
```

## dnscli

```
$ dnscli
Usage: dnscli [flags] command [arguments]

$ dnscli --help
Usage of dnscli:
  -apikey string
    	The request 'x-api-key' (default "default")
  -h string
    	The server host (default "localhost")
  -p int
    	The server port (default 9999)
```
