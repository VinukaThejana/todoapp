// Package handler provides the handlers for the application.
package handler

import (
	"net/http"

	"github.com/bytedance/sonic"
)

func jsonresponse(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	sonic.ConfigDefault.NewEncoder(w).Encode(map[string]string{
		"message": msg,
	})
}
