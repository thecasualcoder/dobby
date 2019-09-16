# dobby

![Dobby GIF](dobby.gif)

dobby is **free** and will serve your orders.

You can start dobby using:

```
dobby server
```

which will start dobby server in port `4444`.

You can ask dobby's

- health `curl dobby:4444/health`
- readiness `curl dobby:4444/readiness`
- version `curl dobby:4444/version`

You can order dobby to

- be healthy

  `PUT /control/health/perfect` which will make `/health` to return 200

- fall sick

  `PUT /control/health/sick` which will make `/health` to return 500

- be ready

  `PUT /control/ready/perfect` which will make `/readiness` to return 200

- not to be ready

  `PUT /control/ready/sick` which will make `/readiness` to return 500

- kill itself

  `PUT /state/crash` which will crash the server

## Run

### Configurations

Available configurations:

| Key               | Value  | Purpose                                   | Default   |
| ----------------- | ------ | ----------------------------------------- | --------- |
| VERSION           | String | To set the version of program             | Build Arg |
| INITIAL_HEALTH    | String | Sets the initial health of the program    | TRUE      |
| INITIAL_READINESS | String | Sets the initial readiness of the program | TRUE      |
| PORT              | Int    | Sets the port of the server               | 4444      |
| BIND_ADDR         | String | Listen address of the process             | 127.0.0.1 |

### Docker

```
docker run thecasualcoder/dobby
```

### Local

```
git clone https://github.com/thecasualcoder/dobby.git && cd dobby
make compile
./out/dobby server
```
