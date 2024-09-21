set dotenv-load

connect:
    usql $(echo $DATABASE_URL)

migrate:
    go run cmd/migrator/main.go

schema file:
    usql $(echo $DATABASE_URL) -f {{ file }}
