# dobby

[![Build Status](https://travis-ci.org/thecasualcoder/dobby.svg?branch=master)](https://travis-ci.org/thecasualcoder/dobby)
[![Go Doc](https://godoc.org/github.com/thecasualcoder/dobby?status.svg)](https://godoc.org/github.com/thecasualcoder/dobby)
[![Go Report Card](https://goreportcard.com/badge/github.com/thecasualcoder/dobby)](https://goreportcard.com/report/github.com/thecasualcoder/dobby)

![Dobby GIF](dobby.gif)

dobby is **free** and will serve your orders.

You can start dobby in Docker using:

```bash
$ docker run -p 4444:4444 thecasualcoder/dobby
```

which will start dobby server in port `4444`.

You can ask dobby's

- health

  ```bash
  curl -i localhost:4444/health
  ```
- readiness
  ```bash
  curl localhost:4444/readiness
  ```

- version
  ```bash
  curl localhost:4444/version
  ```

- metadata about the host
  ```bash
  curl localhost:4444/meta
  ```

You can order dobby to

- Repeat the status code (with requested delay in milliseconds)
  ```bash
  $ curl localhost:4444/return/200 -i
  HTTP/1.1 200 OK
  Date: Sun, 17 May 2020 09:51:13 GMT
  Content-Length: 0

  $ curl localhost:4444/return/401?delay=300 -i
  HTTP/1.1 401 Unauthorized
  Date: Sun, 17 May 2020 09:50:34 GMT
  Content-Length: 0
  ```
- be healthy

  `PUT /control/health/perfect` which will make `/health` to return 200

- fall sick

  `PUT /control/health/sick` which will make `/health` to return 500

- recover health after sometime

  `PUT /control/health/sick?resetInSeconds=2` which will make `/health` to return 500 for 2 seconds

- be ready

  `PUT /control/ready/perfect` which will make `/readiness` to return 200

- not to be ready

  `PUT /control/ready/sick` which will make `/readiness` to return 503

- recover readiness after sometime

  `PUT /control/ready/sick?resetInSeconds=2` which will make `/readiness` to return 503 for 2 seconds

- add load on memory

  `PUT /control/goturbo/memory` which will create a memory spike

- add load on CPU

  `PUT /control/goturbo/cpu` which will create a CPU spike

- kill itself

  `PUT /control/crash` which will crash the server

- call another service

   ```sh
   POST /call -d '{"url": "http://httpbin.org/get", "method": "GET"}' will make a get request to http://httpbin.org/get
   POST /call -d '{"url": "http://httpbin.org/post", "method": "POST", "body": "{key: value}"}' will make a post request to http://httpbin.org/post
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

```bash
$ git clone https://github.com/thecasualcoder/dobby.git && cd dobby
$ make compile
$ ./out/dobby server
```

### Swagger Docs

Swagger docs will be available at: [http://localhost:4444/swagger/index.html](http://localhost:4444/swagger/index.html)

### Contributing

Fork the repo and start contributing

#### Guidelines

- Make sure to run build before raising PR (`make build`)
- Update README.md if necessary
