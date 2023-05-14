# Proxy Web App 
## IOS test task

Built on [Go](https://go.dev/)

## Features

- Hexagon structure
- Routes your requests to another host
- Modifies JSON response by adding `custom_key` to it
- Has adjustable ratelimit


## Installation

Requires [Go](https://go.dev/doc/install) to run application 
Create `.env` file and fill it as it shown in `env.sample`

You will require MongoDB locally or [Cloud solution](https://www.mongodb.com/atlas/database)
Create DB named `core` with collection named `request`

To run application:
```sh
go run server.go
```
To run tests:
```sh
go test -v
```