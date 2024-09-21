package lib

import (
	"net/http"

	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// WrapHandler wraps the handler function with the environment, database, and Redis client.
func WrapHandler(
	h func(http.ResponseWriter, *http.Request, *env.Env, *gorm.DB, *redis.Client),
	e *env.Env,
	db *gorm.DB,
	rdb *redis.Client,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r, e, db, rdb)
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
