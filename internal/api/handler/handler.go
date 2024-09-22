// Package handler provides the handlers for the application.
package handler

import (
	"net/http"

	"github.com/bytedance/sonic"
)

// JSONr is a helper function to send a JSON response
func JSONr(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	sonic.ConfigDefault.NewEncoder(w).Encode(map[string]string{
		"message": msg,
	})
}
