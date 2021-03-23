# dobby

[![Build Status](https://github.com/thecasualcoder/dobby/workflows/Check/badge.svg)](https://github.com/thecasualcoder/dobby/actions?query=workflow%3ACheck+branch%3Amaster)
[![Go Doc](https://godoc.org/github.com/thecasualcoder/dobby?status.svg)](https://godoc.org/github.com/thecasualcoder/dobby)
[![Go Report Card](https://goreportcard.com/badge/github.com/thecasualcoder/dobby)](https://goreportcard.com/report/github.com/thecasualcoder/dobby)

![Dobby GIF](dobby.gif)

## About

dobby is **free** and will serve your orders.

You can start dobby in Docker using:

```shell
$ docker run -p 4444:4444 thecasualcoder/dobby
```

which will start dobby server in port `4444`.

## Features

You can ask dobby any of the following

- [Version](#version)
- [Metadata](#metadata)
- [Health](#health)
    + [About its health](#about-its-health)
    + [To be healthy](#to-be-healthy)
    + [To fall sick](#to-fall-sick)
    + [To recover health after sometime](#to-recover-health-after-sometime)
- [Readiness](#readiness)
    + [About its readiness](#about-its-readiness)
    + [To be ready](#to-be-ready)
    + [To be unready](#to-be-unready)
    + [To recover readiness after sometime](#to-recover-readiness-after-sometime)
- [Disruptions](#disruptions)
    + [Add load on memory](#add-load-on-memory)
    + [Add load on CPU](#add-load-on-cpu)
    + [Kill itself](#kill-itself)
- [Repeat Http Code](#repeat-http-code)
    + [To return a given status code](#to-return-a-given-status-code)
    + [To return a given status code (with requested delay in milliseconds)](#to-return-a-given-status-code-with-requested-delay-in-milliseconds)
- [Call a service](#call-a-service)
    + [To call another service](#to-call-another-service)
- [Configure Proxies](#configure-proxies)
    + [To proxy a call](#to-proxy-a-call)
    + [To delete a configured proxy](#to-delete-a-configured-proxy)

### Version

```shell
$ curl localhost:4444/version
```

### Metadata

```shell
$ curl localhost:4444/meta
```

### Health

Ask dobby

#### About its health
	
```shell
$ curl -i localhost:4444/health
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:32:02 GMT
Content-Length: 16
  
{"healthy":true}
```

#### To be healthy

```shell
$ curl -i -X PUT localhost:4444/control/health/perfect
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:33:02 GMT
Content-Length: 20

{"status":"success"}

$ curl -i localhost:4444/health
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:33:16 GMT
Content-Length: 16

{"healthy":true}
```

#### To fall sick

```shell
$ curl -i -X PUT localhost:4444/control/health/sick
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:33:29 GMT
Content-Length: 20

{"status":"success"}

$ curl -i localhost:4444/health
HTTP/1.1 500 Internal Server Error
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:33:48 GMT
Content-Length: 17

{"healthy":false}
```

#### To recover health after sometime

```shell
$ curl -i -X PUT "localhost:4444/control/health/sick?resetInSeconds=20"
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:34:42 GMT
Content-Length: 20

{"status":"success"}

$ curl -i localhost:4444/health
HTTP/1.1 500 Internal Server Error
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:34:47 GMT
Content-Length: 17

{"healthy":false}

#  After 20seconds
$ curl -i localhost:4444/health
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:35:12 GMT
Content-Length: 16

{"healthy":true}  
```

### Readiness

Ask dobby

#### About its readiness

```shell
$ curl -i localhost:4444/readiness
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:40:30 GMT
Content-Length: 14

{"ready":true}
```

#### To be ready

```shell
$ curl -i -X PUT localhost:4444/control/ready/perfect
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:40:51 GMT
Content-Length: 20

{"status":"success"}

$ curl -i localhost:4444/readiness
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:41:25 GMT
Content-Length: 14

{"ready":true}
```

#### To be unready

```shell
$ curl -i -X PUT localhost:4444/control/ready/sick
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:41:41 GMT
Content-Length: 20

{"status":"success"}

$ curl -i localhost:4444/readiness
HTTP/1.1 503 Service Unavailable
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:41:58 GMT
Content-Length: 15

{"ready":false}
```

#### To recover readiness after sometime

```shell
$ curl -i -X PUT "localhost:4444/control/ready/sick?resetInSeconds=20"
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:42:39 GMT
Content-Length: 20

{"status":"success"}

$ curl -i localhost:4444/readiness
HTTP/1.1 503 Service Unavailable
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:42:43 GMT
Content-Length: 15

{"ready":false}

#  After 20seconds
$ curl -i localhost:4444/readiness
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:43:05 GMT
Content-Length: 14

{"ready":true}
```

### Disruptions

You can also ask dobby to

#### Add load on memory

```shell
# to create a memory spike
$ curl -i -X PUT localhost:4444/control/goturbo/memory
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 12:00:22 GMT
Content-Length: 20

{"status":"success"}
```

#### Add load on CPU

```shell
# to create a cpu spike
$ curl -i -X PUT localhost:4444/control/goturbo/cpu
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 12:01:07 GMT
Content-Length: 20

{"status":"success"}
```

#### Kill itself

```shell
$ curl -i -X PUT localhost:4444/control/crash
# Beware, this is stop the running server
```

### Repeat Http Code

Ask dobby

#### To return a given status code

```shell
$ curl -i localhost:4444/return/200
HTTP/1.1 200 OK
Date: Sun, 17 May 2020 09:51:13 GMT
Content-Length: 0
```

#### To return a given status code (with requested delay in milliseconds)

```shell
$ curl -i localhost:4444/return/401?delay=300
HTTP/1.1 401 Unauthorized
Date: Sun, 17 May 2020 09:50:34 GMT
Content-Length: 0
```

### Call a service

#### To call another service

```shell
# To make a get request to http://httpbin.org/get
$ curl -i localhost:4444/call -d '{"url": "http://httpbin.org/get", "method": "GET"}'
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:46:25 GMT
Content-Length: 240

{"args":{},"headers":{},"url":"http://httpbin.org/get"}

# To make a post request to http://httpbin.org/post
$ curl -i localhost:4444/call -d '{"url": "http://httpbin.org/post", "method": "POST", "body": "{key: value}"}'
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:49:09 GMT
Content-Length: 311

{"args":{},"data":"\"{key: value}\"","files":{},"form":{},"headers":{},"json":"{key: value}","url":"http://httpbin.org/post"}
```

### Configure Proxies
  
#### To proxy a call

```shell
$ curl -i localhost:4444/proxy -d '{"path":"/time","method": "GET", "proxy": {"url":"http://worldtimeapi.org/api/timezone/asia/kolkata","method":"GET"}}'
HTTP/1.1 201 Created
Date: Tue, 16 Mar 2021 11:51:32 GMT
Content-Length: 0

$ curl -i localhost:4444/time
# makes a call to http://worldtimeapi.org/api/timezone/asia/kolkata
```

#### To delete a configured proxy

```shell
$ curl -i -X DELETE localhost:4444/proxy -d '{"path":"/time","method": "GET", "proxy": {"url":"http://worldtimeapi.org/api/timezone/asia/kolkata","method":"GET"}}'
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Mar 2021 11:55:22 GMT
Content-Length: 50

{"result":"deleted the proxy config successfully"}

$ curl -i localhost:4444/time
HTTP/1.1 404 Not Found
Content-Type: text/plain
Date: Tue, 16 Mar 2021 11:55:46 GMT
Content-Length: 18

404 page not found
```

## Run

### Configurations

Available configurations:

| Key               | Value  | Purpose                                                    | Default   |
| ----------------- | ------ | ---------------------------------------------------------- | --------- |
| VERSION           | String | To set the version of program                              | Build Arg |
| INITIAL_DELAY     | Int    | Sets the initial delay to start the server (in seconds)    | 0         |
| INITIAL_HEALTH    | String | Sets the initial health of the program                     | TRUE      |
| INITIAL_READINESS | String | Sets the initial readiness of the program                  | TRUE      |
| PORT              | Int    | Sets the port of the server                                | 4444      |
| BIND_ADDR         | String | Listen address of the process                              | 127.0.0.1 |

### Run in local

```shell
$ git clone https://github.com/thecasualcoder/dobby.git && cd dobby
$ make compile
$ ./out/dobby server
```

### Swagger Docs

Swagger docs will be available at: [http://localhost:4444/swagger/index.html](http://localhost:4444/swagger/index.html)

## Contributing

Fork the repo and start contributing

#### Guidelines

- Make sure to run build before raising PR (`make build`)
- Make sure to generate and check in swagger docs if any added (`make swagger-docs`)
- Update README.md if necessary
