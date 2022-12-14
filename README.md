## Short link service

Golang service with UI for processing short links like Google, but more pleasant for your company domain

## Features

- Web UI admin panel with Basic auth
- VanillaJS
- Interface based arch
- Adopted for PostgreSQL

## Installation

Install Link Shortnet

- Install [go 1.18](https://go.dev/doc/install) (latest)
- Clone repo
- `go run main.go`


## Build

docker build --tag shortlink -t shortlink:multistage .

## Run

docker run --env-file .env_test -p 8080:8080 shortlink:multistage


## License

[GPL v3](https://www.gnu.org/licenses/gpl-3.0.en.html)
