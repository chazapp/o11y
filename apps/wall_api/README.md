# Wall-API

[![WallAPI](https://github.com/chazapp/o11y/actions/workflows/wall_api_tests.yaml/badge.svg)](https://github.com/chazapp/o11y/actions/workflows/wall_api_tests.yaml)
[![Coverage](https://codecov.io/gh/chazapp/o11y/graph/badge.svg?token=FIAGTCSSD1&flag=wall-api)](https://codecov.io/gh/chazapp/o11y)  

## Features

- A CRUD API for Messages
- ORM Layer to PostgreSQL via `go-orm`
- A WebSocket server that emits `Message` objects on reception of `POST /message`


## How to run

Have a Postgres instance running via Docker:

```bash
o11y/apps/wall_api$ docker run -e POSTGRES_DB=wallapi -e POSTGRES_USER=user -e POSTGRES_PASSWORD=foobar -p 5432:5432 -d  postgres:16
```

Build and start the API

```bash
o11y/apps/wall_api$ go build -o wallapi
o11y/apps/wall_api$ ./wallapi
NAME:
   wall-api - An API for the Wall Application

USAGE:
   wall-api [global options] command [command options] [arguments...]

COMMANDS:
   run      Run the API server
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

o11y/apps/wall_api$ ./wallapi run --dbHost 127.0.0.1 --dbUser user --dbPassword foobar --port 8080 --dbName wallapi
```

Run the test suite:

```bash
$ go test ./... -cover
$ go test ./... -cover
ok      github.com/chazapp/o11/apps/wall_api    0.450s  coverage: 26.4% of statements
?       github.com/chazapp/o11/apps/wall_api/api        [no test files]
?       github.com/chazapp/o11/apps/wall_api/metrics    [no test files]
?       github.com/chazapp/o11/apps/wall_api/models     [no test files]
?       github.com/chazapp/o11/apps/wall_api/ws [no test files]
```