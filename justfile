set dotenv-load := true

default: dev

dev:
    air

PG_CONNSTR := env("POSTGRES_CONNSTR")

migrate:
    goose postgres {{ PG_CONNSTR }} up -dir ./store/migrations

migrate-one:
    goose postgres {{ PG_CONNSTR }} up-by-one -dir ./store/migrations

migrate-down:
    goose postgres {{ PG_CONNSTR }} down -dir ./store/migrations

hanged:
    pkill district-dev
