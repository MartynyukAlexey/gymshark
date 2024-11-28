package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/MartynyukAlexey/gymshark/internal/service/auth"
)

func HandleRefresh(svc *auth.Service, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			serveError(w, "no refresh token", http.StatusUnauthorized)
			return
		}

		refreshResp, err := svc.Refresh(&auth.RefreshReq{
			RefreshToken: cookie.Value,
			Email:        "MartynyukAlexey05@gmail.com",
		})

		if err != nil {
			switch err {
			case auth.ErrInvalidRefreshToken:
				serveError(w, err.Error(), http.StatusUnauthorized)
			case auth.ErrUserNotFound:
				serveError(w, err.Error(), http.StatusNotFound)
			default:
				serveError(w, "internal error", http.StatusInternalServerError)
			}

			return
		}

		cookieAccessToken := &http.Cookie{
			Name:     "access_token",
			Value:    refreshResp.AccessToken,
			HttpOnly: true,
			MaxAge:   int(svc.Cfg.AccessTokenTTL.Seconds()),
		}

		cookieRefreshToken := &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshResp.RefreshToken,
			HttpOnly: true,
			MaxAge:   int(svc.Cfg.RefreshTokenTTL.Seconds()),
		}

		http.SetCookie(w, cookieAccessToken)
		http.SetCookie(w, cookieRefreshToken)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(
			struct {
				Status string `json:"status"`
				Msg    string `json:"message"`
			}{
				Status: "ok",
				Msg:    "successful login",
			},
		)

	}
}
