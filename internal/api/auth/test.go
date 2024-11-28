package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/MartynyukAlexey/gymshark/internal/service/auth"
)

func HandleTest(svc *auth.Service, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(
			struct {
				Status  string `json:"status"`
				Message string `json:"message"`
			}{
				Status:  "ok",
				Message: "success",
			})
	}
}
