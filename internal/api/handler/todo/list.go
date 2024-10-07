// Package todo : This package is for getting all the todo for a given user
package todo

import (
	"net/http"

	"github.com/VinukaThejana/todoapp/internal/api/grpc"
	"github.com/VinukaThejana/todoapp/internal/api/middleware"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/pkg/todo"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// List: This function is for getting all the todos under a given user
func List(
	w http.ResponseWriter,
	r *http.Request,
	tcm *grpc.TodoClientManager,
	e *env.Env,
	db *gorm.DB,
	rdb *redis.Client,
) {
	userID := r.Context().Value(middleware.UserID).(string)

	res, err := tcm.Client().List(r.Context(), &todo.ListRequest{
		UserId: userID,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to update the todo")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		sonic.ConfigDefault.NewEncoder(w).Encode(res.Todos)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	sonic.ConfigDefault.NewEncoder(w).Encode(res.Todos)
	return
}
