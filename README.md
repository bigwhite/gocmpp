[![Build Status](https://travis-ci.org/bigwhite/gocmpp.svg?branch=master)](https://travis-ci.org/bigwhite/gocmpp)
[![Coverage Status](https://coveralls.io/repos/github/bigwhite/gocmpp/badge.svg?branch=master)](https://coveralls.io/github/bigwhite/gocmpp?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/bigwhite/gocmpp)](https://goreportcard.com/report/github.com/bigwhite/gocmpp)
[![GoDoc](http://img.shields.io/badge/GoDoc-Reference-blue.svg)](https://godoc.org/github.com/bigwhite/gocmpp)

# gocmpp
An implementation of China Mobile Peer to Peer(cmpp) protocol in golang for both client and server sides.

The protocol versions below are covered：
* [cmpp2.1](http://pan.baidu.com/s/13E0Q6)
* [cmpp3.0](http://pan.baidu.com/s/1o61obA6)

## Overview
gocmpp is a library that implements China Mobile Peer to Peer(cmpp) protocol. You can use the lib to implement any app, tool or system that use cmpp protocol for both client and server side.

gocmpp is portable well. It can be used on linux, darwin or even windows.

gocmpp has covered cmpp2.x and cmpp3.x versions. It has supported the connect, submit, deliver, fwd, active test, and terminate packets of cmpp protocol. But other less use packets like cmpp query, cancel and route have not been supported. And they are not in the roadmap at all.

## QuickStart

### 1. Download gocmpp
```
$go get github.com/bigwhite/gocmpp
```

### 2. Build gocmpp

gocmpp's build use the vendor mechanism introduced in go 1.5, so we need go 1.5 or later version to build it.

```
$export GO15VENDOREXPERIMENT="1"
$cd $GOPATH/src/github.com/bigwhite/gocmpp
$make
$make
go build -o examples/server/server examples/server/server.go
go build -o examples/client/client examples/client/client.go
go build
go build ./utils
go test
PASS
ok  	github.com/bigwhite/gocmpp	0.009s
go test ./utils
ok  	github.com/bigwhite/gocmpp/utils	0.008s
```

### 3. Run the examples

run the two programs below in order:
 * ./examples/server/server
 * ./examples/client/client

you would get the output like below:

server:
```
cmppserver: 2015/11/19 16:28:50 accept a connection from 127.0.0.1:49847
cmppserver: 2015/11/19 16:28:50 receive a cmpp30 connect request from 127.0.0.1:49847[0]
cmppserver: 2015/11/19 16:28:50 handleLogin: 900001 login ok
cmppserver: 2015/11/19 16:28:55 receive a cmpp active response from 127.0.0.1:49847[0]
cmppserver: 2015/11/19 16:29:00 receive a cmpp active response from 127.0.0.1:49847[1]
cmppserver: 2015/11/19 16:29:00 receive a cmpp30 submit request from 127.0.0.1:49847[1]
cmppserver: 2015/11/19 16:29:00 handleSubmit: handle submit from 900001 ok! msgid[12878564852733378560], srcId[900001], destTerminalId[13500002696]
cmppserver: 2015/11/19 16:29:05 receive a cmpp active response from 127.0.0.1:49847[2]
cmppserver: 2015/11/19 16:29:05 receive a cmpp30 submit request from 127.0.0.1:49847[2]
cmppserver: 2015/11/19 16:29:05 handleSubmit: handle submit from 900001 ok! msgid[12878564852733378560], srcId[900001], destTerminalId[13500002696]
cmppserver: 2015/11/19 16:29:10 receive a cmpp active response from 127.0.0.1:49847[3]
cmppserver: 2015/11/19 16:29:10 receive a cmpp30 submit request from 127.0.0.1:49847[3]
cmppserver: 2015/11/19 16:29:10 handleSubmit: handle submit from 900001 ok! msgid[12878564852733378560], srcId[900001], destTerminalId[13500002696]
cmppserver: 2015/11/19 16:29:13 close connection with 127.0.0.1:49847!
```

client:

```
2015/11/19 16:28:50 client connect and auth ok
2015/11/19 16:28:55 receive a cmpp active request: &{0}
2015/11/19 16:29:00 receive a cmpp active request: &{1}
2015/11/19 16:29:00 send a cmpp3 submit request
2015/11/19 16:29:00 receive a cmpp3 submit response: &{12878564852733378560 0 1}
2015/11/19 16:29:05 receive a cmpp active request: &{2}
2015/11/19 16:29:05 send a cmpp3 submit request
2015/11/19 16:29:05 receive a cmpp3 submit response: &{12878564852733378560 0 2}
2015/11/19 16:29:10 receive a cmpp active request: &{3}
2015/11/19 16:29:10 send a cmpp3 submit request
2015/11/19 16:29:10 receive a cmpp3 submit response: &{12878564852733378560 0 3}
```
### 4.How to use the library
The best start entry is the code in folder '[examples](https://github.com/bigwhite/gocmpp/tree/master/examples)'. 
 * client.go tells you how to use gocmpp in client side.
 * server.go presents you how to construct a cmpp server in gocmpp.
 
## Who use gocmpp
* [smspush platform](https://51smspush.com)
