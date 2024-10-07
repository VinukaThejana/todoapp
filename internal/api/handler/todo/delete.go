// Package todo : This package is for deleting a given event
package todo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/VinukaThejana/todoapp/internal/api/grpc"
	"github.com/VinukaThejana/todoapp/internal/api/handler"
	"github.com/VinukaThejana/todoapp/internal/api/middleware"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/pkg/todo"
	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Delete: This function is for deleting a given todo
func Delete(
	w http.ResponseWriter,
	r *http.Request,
	tcm *grpc.TodoClientManager,
	e *env.Env,
	db *gorm.DB,
	rdb *redis.Client,
) {
	const (
		maxRequestBodySize = 1 << 14
	)

	type body struct {
		ID uint `json:"id" validate:"required"`
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	defer r.Body.Close()

	var reqBody body

	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Error().Err(err)
		handler.JSONr(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	validate := validator.New()
	err = validate.Struct(reqBody)
	if err != nil {
		log.Error().Err(err).Msg("validation failed")

		validationErrs := err.(validator.ValidationErrors)
		handler.JSONr(w, http.StatusBadRequest, fmt.Sprintf("Please provide a valid %s", strings.ToLower(validationErrs[0].Field())))
		return
	}

	userID := r.Context().Value(middleware.UserID).(string)

	_, err = tcm.Client().Delete(r.Context(), &todo.DeleteRequest{
		Id:     fmt.Sprint(reqBody.ID),
		UserId: userID,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to delete the todo")
		st, ok := status.FromError(err)
		if !ok {
			handler.JSONr(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		switch st.Code() {
		case codes.Unauthenticated:
			handler.JSONr(w, http.StatusUnauthorized, "Unauthenticated")
			return
		case codes.NotFound:
			handler.JSONr(w, http.StatusNotFound, "Todo not found")
			return
		default:
			handler.JSONr(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	handler.JSONr(w, http.StatusCreated, "Todo updated successfully")
}
