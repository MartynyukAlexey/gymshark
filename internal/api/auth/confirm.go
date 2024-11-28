package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/MartynyukAlexey/gymshark/internal/service/auth"
)

func HandleConfirmation(svc *auth.Service, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req auth.ConfirmReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			serveError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := svc.Confirm(&req); err != nil {
			switch err {
			case auth.ErrInvalidCode:
				serveError(w, err.Error(), http.StatusForbidden)
			case auth.ErrUserNotFound:
				serveError(w, err.Error(), http.StatusNotFound)
			case auth.ErrUserAlreadyConfirmed:
				serveError(w, err.Error(), http.StatusConflict)
			case auth.ErrCodeExpired:
				serveError(w, err.Error(), http.StatusGone)
			default:
				serveError(w, "internal error", http.StatusInternalServerError)
			}

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		json.NewEncoder(w).Encode(
			struct {
				Status  string `json:"status"`
				Message string `json:"message"`
			}{
				Status:  "ok",
				Message: "account was confirmed",
			},
		)
	}
}
