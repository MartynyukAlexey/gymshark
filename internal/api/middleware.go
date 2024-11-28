package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/MartynyukAlexey/gymshark/internal/service/auth"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
)

type middlewareEnv struct {
	svc    *auth.Service
	logger *slog.Logger
}

func (env *middlewareEnv) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			serveError(w, "no access token", http.StatusUnauthorized)
			return
		}

		userID, err := env.svc.Authorize(cookie.Value)
		if err != nil {
			serveError(w, "invalid access token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
