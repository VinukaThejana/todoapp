package auth

import (
	"net/http"
	"time"

	"github.com/VinukaThejana/todoapp/internal/api/grpc"
	"github.com/VinukaThejana/todoapp/internal/api/handler"
	"github.com/VinukaThejana/todoapp/internal/api/middleware"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/internal/enums"
	"github.com/VinukaThejana/todoapp/pkg/auth"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Logout logs out the user.
func Logout(
	w http.ResponseWriter,
	r *http.Request,
	acm *grpc.AuthClientManager,
	e *env.Env,
	db *gorm.DB,
	rdb *redis.Client,
) {
	refreshToken := r.Context().Value(middleware.RefreshToken).(string)
	_, err := acm.Client().Logout(r.Context(), &auth.LogoutRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Error().Err(err)
			handler.JSONr(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		switch st.Code() {
		case codes.Unauthenticated:
			handler.JSONr(w, http.StatusUnauthorized, "Not authorized to perform this action")
			return
		default:
			handler.JSONr(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "todoapp_refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().UTC().Add(-24 * time.Hour),
		Secure:   e.Environ == string(enums.Prd),
		Domain:   e.Domain,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "todoapp_session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		Expires:  time.Now().UTC().Add(-24 * time.Hour),
		Secure:   e.Environ == string(enums.Prd),
		Domain:   e.Domain,
	})

	handler.JSONr(w, http.StatusOK, "Successfully logged out")
}
