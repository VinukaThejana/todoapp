// Package middleware provides the middleware for the application.
package middleware

import (
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
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
