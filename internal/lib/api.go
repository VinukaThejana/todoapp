package lib

import (
	"net/http"

	"github.com/VinukaThejana/todoapp/internal/api/grpc"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// WrapHandlerWAuthClient wraps the handler function with the environment, database, Redis client, and authentication service client.
func WrapHandlerWAuthClient(
	h func(http.ResponseWriter, *http.Request, *grpc.AuthClientManager, *env.Env, *gorm.DB, *redis.Client),
	acm *grpc.AuthClientManager,
	e *env.Env,
	db *gorm.DB,
	rdb *redis.Client,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r, acm, e, db, rdb)
	}
}

// WrapHandlerWTodoClient wraps the handler function with the environment, database, Redis client, and todo service client.
func WrapHandlerWTodoClient(
	h func(http.ResponseWriter, *http.Request, *grpc.TodoClientManager, *env.Env, *gorm.DB, *redis.Client),
	tcm *grpc.TodoClientManager,
	e *env.Env,
	db *gorm.DB,
	rdb *redis.Client,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r, tcm, e, db, rdb)
	}
}

// WrapHandler wraps the handler function with the environment, database, Redis client, and both the authentication and todo service clients.
func WrapHandler(
	h func(http.ResponseWriter, *http.Request, *grpc.AuthClientManager, *grpc.TodoClientManager, *env.Env, *gorm.DB, *redis.Client),
	acm *grpc.AuthClientManager,
	tcm *grpc.TodoClientManager,
	e *env.Env,
	db *gorm.DB,
	rdb *redis.Client,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r, acm, tcm, e, db, rdb)
	}
}

// WrapMiddleware wraps the middleware function with the environment, database, and Redis client.
func WrapMiddleware(
	m func(http.Handler, *env.Env, *gorm.DB, *redis.Client) http.Handler,
	e *env.Env,
	db *gorm.DB,
	rdb *redis.Client,
) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return m(h, e, db, rdb)
	}
}

// WrapMiddlewareWAuth wraps the middleware function with the environment, database, Redis client, and authentication service client.
func WrapMiddlewareWAuth(
	m func(http.Handler, *grpc.AuthClientManager, *env.Env, *gorm.DB, *redis.Client) http.Handler,
	acm *grpc.AuthClientManager,
	e *env.Env,
	db *gorm.DB,
	rdb *redis.Client,
) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return m(h, acm, e, db, rdb)
	}
}
