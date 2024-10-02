// Package middleware provides the middleware for the application.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/VinukaThejana/todoapp/internal/api/grpc"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/pkg/auth"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

const (
	RefreshToken = "todoapp_refresh_token"
	UserID       = "user_id"
)

// ContentJSON is a middleware that checks if the content type is application/json
func ContentJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			log.Error().
				Msgf(
					"Content-Type : %s\tinvalid content type provided",
					contentType,
				)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			sonic.ConfigDefault.NewEncoder(w).Encode(
				map[string]interface{}{
					"message": "only conent type of application/json is allowed",
				})
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RefreshTokenPresent is a middleware that checks if the refresh token is present in the request
func RefreshTokenPresent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshTokenC, err := r.Cookie("todoapp_refresh_token")
		if err != nil {
			log.Error().Msg("refresh token not present")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), RefreshToken, refreshTokenC.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Auth is a middleware that validates the access token and assigns the user id of the requesting
// user to the context if the access token is valid.
func Auth(next http.Handler, acm *grpc.AuthClientManager, e *env.Env, db *gorm.DB, rdb *redis.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			log.Error().Msg("auth token not present")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		authToken := strings.TrimPrefix(header, "Bearer ")

		res, err := acm.Client().Validate(r.Context(), &auth.ValidateRequest{
			AccessToken: authToken,
		})
		if err != nil || !res.IsValid || !res.Success {
			st, ok := status.FromError(err)
			if !ok {
				log.Error().Err(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			switch st.Code() {
			case codes.Unauthenticated:
				w.WriteHeader(http.StatusUnauthorized)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		ctx := context.WithValue(r.Context(), UserID, res.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
