// Package router provides the routing for the application.
package router

import (
	"github.com/VinukaThejana/todoapp/internal/api/grpc"
	"github.com/VinukaThejana/todoapp/internal/api/handler/auth"
	m "github.com/VinukaThejana/todoapp/internal/api/middleware"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/internal/lib"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Init(
	acm *grpc.AuthClientManager,
	tcm *grpc.TodoClientManager,
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
			r.Post("/register", lib.WrapHandlerWAuthClient(
				auth.Register,
				acm, e, db, rdb,
			))
			r.Post("/login", lib.WrapHandlerWAuthClient(
				auth.Login,
				acm, e, db, rdb,
			))
		})

		r.Group(func(r chi.Router) {
			r.Use(m.RefreshTokenPresent)
			r.Patch("/refresh", lib.WrapHandlerWAuthClient(
				auth.Login,
				acm, e, db, rdb,
			))
			r.Delete("/logout", lib.WrapHandlerWAuthClient(
				auth.Logout,
				acm, e, db, rdb,
			))
		})
	})

	return r
}
