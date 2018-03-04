# Mika
 一个简单快速的安全代理:rocket:

## 特性

1. 支持本地 socks5 代理 TCP 流量
2. 支持本地 HTTP 代理
3. 支持无特征流量加密
4. 支持通过 HTTP 流量混淆

## 目标

1. 更少的协议特征
2. 足够安全保护你的因特网流量
3. 传输速度够快
4. 任何新的特性或者功能增强不能破坏以上目标

## Mika 协议说明

见 [Mika Protocol Spec](https://github.com/sakeven/mika/wiki/Mika-Protocol-Spec)

## 配置
### 客户端
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
### 服务端

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
### 通过配置文件配置
见 [Configuration via Config File](https://github.com/sakeven/mika/wiki/Configuration-via-Config-File)

## 构建

```
./build.sh
```

两个二进制程序 `client` 和 `server` 将会被安装在代码根目录的 `bin/` 下。

## LICENSE

In MIT LICENSE
