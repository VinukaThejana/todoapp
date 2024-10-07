set dotenv-load

mod db
mod redis

auth:
  go run cmd/auth/main.go

todo:
  go run cmd/todo/main.go

run:
  go run cmd/api/main.go

docker-compose:
  docker compose -f deployments/docker-compose.yml up
