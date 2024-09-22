package auth

import (
	"net/http"

	"github.com/VinukaThejana/todoapp/internal/api/grpc"
	"github.com/VinukaThejana/todoapp/internal/api/handler"
	"github.com/VinukaThejana/todoapp/internal/api/middleware"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/pkg/auth"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Refresh is the handler for the refresh endpoint.
func Refresh(
	w http.ResponseWriter,
	r *http.Request,
	acm *grpc.AuthClientManager,
	e *env.Env,
	db *gorm.DB,
	rdb *redis.Client,
) {
	refreshToken := r.Context().Value(middleware.RefreshToken).(string)
	resp, err := acm.Client().Refresh(r.Context(), &auth.RefreshRequest{
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
			handler.JSONr(w, http.StatusUnauthorized, "Invalid password")
			return
		default:
			handler.JSONr(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	w.Header().Add("X-New-Access-Token", resp.AccessToken)
	handler.JSONr(w, http.StatusOK, "Refreshed")
	return
}
