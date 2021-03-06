# Mika
A Faster Secure Proxy :rocket:

[中文指南](README-zh.md)

## Feature

1. Support proxy TCP data by socks5 at local.
2. Support HTTP/HTTPS proxy.
3. Transfer speed is super fast.
4. Support HTTP obfs.

## Goals

1. Less protocol characteristics.
2. Enough security to protect your Internet traffic.
3. Transfer speed should be super fast.
4. Any enhancement or feature shoudn't break above goals.

## Mika Protocol Spec

See [Mika Protocol Spec](https://github.com/sakeven/mika/wiki/Mika-Protocol-Spec)

## Configuration
### Client
```
Usage of client:
  --help
    	print usage
  -b string
    	local binding address (default "127.0.0.1")
  -c string
    	path to config file
  -k string
    	password (default "password")
  -l int
    	local port (default 1080)
  -m string
    	encryption method (default "aes-256-cfb")
  -p int
    	server port (default 8388)
  -s string
    	server address
  -t int
    	timeout in seconds (default 300)
```
### Server

```
Usage of server:
  --help
    	print usage
  -c string
    	path to config file
  -k string
    	password (default "password")
  -m string
    	encryption method (default "aes-256-cfb")
  -p int
    	server port (default 8388)
  -s string
    	server address
  -t int
    	timeout in seconds (default 300)
```
### Configuration via Config File
See [Configuration via Config File](https://github.com/sakeven/mika/wiki/Configuration-via-Config-File)

## Build

```
./build.sh
```

Two binaries `client` and `server` will be installed at `bin/`.

## LICENSE

In MIT LICENSE
