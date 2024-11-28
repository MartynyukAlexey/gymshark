package auth

import (
	"encoding/json"
	"net/http"
)

func serveError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	response := map[string]string{
		"status":  "error",
		"message": msg,
	}

	json.NewEncoder(w).Encode(response)
}
