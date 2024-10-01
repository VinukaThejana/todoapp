package auth

import (
	"net/http"
	"time"

	"github.com/VinukaThejana/todoapp/internal/api/grpc"
	"github.com/VinukaThejana/todoapp/internal/api/handler"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/internal/enums"
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

// Login handler is used for logging in a user
func Login(
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
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password" validate:"required,min=8,max=100,password"`
		Validate string `validate:"validate_login"`
	}

	validate := validator.New()

	validate.RegisterValidation("validate_login", lib.ValidateEmailOrUsername)
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
		log.Error().Err(err)
		handler.JSONr(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var resp *auth.LoginResponse
	if reqBody.Username != "" {
		resp, err = acm.Client().Login(r.Context(), &auth.LoginRequest{
			Login: &auth.LoginRequest_Username{
				Username: reqBody.Username,
			},
			Password: reqBody.Password,
		})
	} else {
		resp, err = acm.Client().Login(r.Context(), &auth.LoginRequest{
			Login: &auth.LoginRequest_Email{
				Email: reqBody.Email,
			},
			Password: reqBody.Password,
		})
	}
	if err != nil {
		st, ok := status.FromError(err)
		log.Error().Err(err).Msg("there is an an error somewhere")
		if !ok {
			handler.JSONr(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		switch st.Code() {
		case codes.InvalidArgument:
			log.Error().Err(err)
			handler.JSONr(w, http.StatusBadRequest, "Invalid request body")
			return
		case codes.NotFound:
			handler.JSONr(w, http.StatusNotFound, "User not found")
			return
		case codes.Unauthenticated:
			handler.JSONr(w, http.StatusUnauthorized, "Invalid password")
			return
		default:
			handler.JSONr(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	w.Header().Add("X-New-Access-Token", resp.TokenSet.AccessToken)
	http.SetCookie(w, &http.Cookie{
		Name:     "todoapp_refresh_token",
		Value:    resp.TokenSet.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(e.RefreshTokenExpiresIn.Seconds()),
		Expires:  time.Now().UTC().Add(e.RefreshTokenExpiresIn),
		Secure:   e.Environ == string(enums.Prd),
		Domain:   e.Domain,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "todoapp_session_token",
		Value:    resp.TokenSet.SessionToken,
		Path:     "/",
		HttpOnly: false,
		MaxAge:   int(e.RefreshTokenExpiresIn.Seconds()),
		Secure:   e.Environ == string(enums.Prd),
		Domain:   e.Domain,
	})

	handler.JSONr(w, http.StatusOK, "Login successful")
}
