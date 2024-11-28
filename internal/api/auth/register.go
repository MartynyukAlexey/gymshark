package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/MartynyukAlexey/gymshark/internal/service/auth"
)

func HandleRegistration(svc *auth.Service, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req auth.RegisterReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			serveError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		authID, err := svc.Register(&req)
		if err != nil {
			switch err {
			case auth.ErrInvalidEmail, auth.ErrInvalidName, auth.ErrInvalidPassword:
				serveError(w, err.Error(), http.StatusBadRequest)
			case auth.ErrUserAlreadyExists:
				serveError(w, err.Error(), http.StatusConflict)
			default:
				serveError(w, "internal error", http.StatusInternalServerError)
			}

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		json.NewEncoder(w).Encode(
			struct {
				Status  string    `json:"status"`
				Message string    `json:"message"`
				UserID  uuid.UUID `json:"user_id"`
			}{
				Status:  "ok",
				Message: "email sent (pending activation)",
				UserID:  authID,
			})
	}
}
