# Wall-API

[![WallAPI](https://github.com/chazapp/o11y/actions/workflows/wall_api.yaml/badge.svg)](https://github.com/chazapp/o11y/actions/workflows/wall_api.yaml)  

## Features

- A CRUD REST API for the `Message` object
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
