// Package middleware provides the middleware for the application.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
)

const (
	RefreshToken = "todoapp_refresh_token"
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

// AuthTokenPresent is a middleware that checks if the auth token is present in the request
func AuthTokenPresent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			log.Error().Msg("auth token not present")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "todoapp_access_token", authToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
