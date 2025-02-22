package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/VinukaThejana/todoapp/internal/api/grpc"
	"github.com/VinukaThejana/todoapp/internal/api/handler"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/internal/lib"
	"github.com/VinukaThejana/todoapp/pkg/auth"
	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Register handler is used for regsitering a new user
func Register(
	w http.ResponseWriter,
	r *http.Request,
	acm *grpc.AuthClientManager,
	e *env.Env,
	db *gorm.DB,
	rdb *redis.Client,
) {
	const (
		maxRequestBodySize = 1 << 8
	)

	type body struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,alphanum,min=4,max=15,username"`
		Name     string `json:"name" validate:"required,min=4,max=30"`
		Password string `json:"password" validate:"required,min=8,max=100,password"`
	}

	validate := validator.New()

	validate.RegisterValidation("username", lib.ValidateUsername)
	validate.RegisterValidation("password", lib.ValidatePassword)

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	defer r.Body.Close()

	var reqBody body

	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Error().Err(err)
		handler.JSONr(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = validate.Struct(reqBody)
	if err != nil {
		log.Error().Err(err).Msg("validation failed")

		validationErrs := err.(validator.ValidationErrors)
		handler.JSONr(w, http.StatusBadRequest, fmt.Sprintf("Please provide a valid %s", strings.ToLower(validationErrs[0].Field())))
		return
	}

	_, err = acm.Client().Register(r.Context(), &auth.RegisterRequest{
		Email:    reqBody.Email,
		Username: reqBody.Username,
		Name:     reqBody.Name,
		Password: reqBody.Password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error().Err(err)
			handler.JSONr(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		switch st.Code() {
		case codes.AlreadyExists:
			handler.JSONr(w, http.StatusConflict, "User already exists")
			return
		default:
			handler.JSONr(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	handler.JSONr(w, http.StatusCreated, "User created")
}
