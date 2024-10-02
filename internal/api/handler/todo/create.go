// Package todo : This package is for creating a new todo item
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
	"gorm.io/gorm"
)

// Create : This function is for creating a new todo item
func Create(
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
		Title       string `json:"title" validate:"required,min=4,max=30"`
		Description string `json:"description" validate:"required,min=4,max=200"`
		Content     string `json:"content" validate:"required,min=4,max=1000"`
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

	_, err = tcm.Client().Create(r.Context(), &todo.CreateRequest{
		Title:       reqBody.Title,
		Description: reqBody.Description,
		Content:     reqBody.Content,
		UserId:      userID,
	})
	if err != nil {
		log.Error().Err(err)
		handler.JSONr(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	handler.JSONr(w, http.StatusCreated, "Todo created")
}
