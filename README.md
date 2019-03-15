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
- version `curl dobby:4444/version`

You can order dobby to

- be healthy
    
    `PUT /state/healthy` which will make `/health` to return 200
- fall sick

    `PUT /state/sick` which will make `/health` to return 500
- kill itself

    `PUT /state/crash` which will crash the server

## Run

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