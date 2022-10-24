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

// Запуск сервиса в контейнере для prod
docker run --rm -d --name passport-verifier \
-e APP_DB_DRIVER=pgx \
-e APP_DB_DSN=postgres://username:password@127.0.0.1:5432/dbname?sslmode=disable \
-e APP_LOG_LEVEL=2 \
-e APP_PORT=8444 \
-e APP_PASSPORTS_URL="https://проверки.гувм.мвд.рф/upload/expired-passports/list_of_expired_passports.csv.bz2" \
-p 8444:8444 passport-verifier:1.0


## License

[GPL v3](https://www.gnu.org/licenses/gpl-3.0.en.html)
