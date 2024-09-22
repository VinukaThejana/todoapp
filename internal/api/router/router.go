// Package router provides the routing for the application.
package router

import (
	"github.com/VinukaThejana/todoapp/internal/api/grpc"
	m "github.com/VinukaThejana/todoapp/internal/api/middleware"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Init(
	acm *grpc.AuthClientManager,
	e *env.Env,
	db *gorm.DB,
	rdb *redis.Client,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/auth", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(m.ContentJSON)
		})

		r.Group(func(r chi.Router) {
		})
	})

	return r
}
