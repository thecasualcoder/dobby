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

### Contributing

Fork the repo and start contributing

#### Guidelines

- Make sure to run build before raising PR (`make build`)
- Update README.md if necessary
