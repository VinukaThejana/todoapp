// Package todo : This package is for getting a todo with the given id
package todo

import (
	"net/http"

	"github.com/VinukaThejana/todoapp/internal/api/grpc"
	"github.com/VinukaThejana/todoapp/internal/api/middleware"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/pkg/todo"
	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Get : This function is for getting a todo with the given id
func Get(
	w http.ResponseWriter,
	r *http.Request,
	tcm *grpc.TodoClientManager,
	e *env.Env,
	db *gorm.DB,
	rdb *redis.Client,
) {
	todoID := chi.URLParam(r, "id")
	if todoID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		sonic.ConfigDefault.NewEncoder(w).Encode(map[string]any{})
		return
	}

	userID := r.Context().Value(middleware.UserID).(string)

	res, err := tcm.Client().Get(r.Context(), &todo.GetRequest{
		Id:     todoID,
		UserId: userID,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to get the todo")
		st, ok := status.FromError(err)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			sonic.ConfigDefault.NewEncoder(w).Encode(map[string]any{})
			return
		}

		switch st.Code() {
		case codes.NotFound:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			sonic.ConfigDefault.NewEncoder(w).Encode(map[string]any{})
			return
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			sonic.ConfigDefault.NewEncoder(w).Encode(map[string]any{})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	sonic.ConfigDefault.NewEncoder(w).Encode(res.Todo)

}
