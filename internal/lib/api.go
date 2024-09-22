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
