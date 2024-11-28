package auth

import (
	"encoding/json"
	"net/http"
)

func serveError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	json.NewEncoder(w).Encode(
		struct {
			Status string `json:"status"`
			Msg    string `json:"message"`
		}{
			Status: "error",
			Msg:    msg,
		},
	)
}
