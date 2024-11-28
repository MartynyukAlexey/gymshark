package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/MartynyukAlexey/gymshark/internal/service/auth"
)

func HandleLogin(svc *auth.Service, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req auth.LoginReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			serveError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		loginResp, err := svc.Login(&req)
		if err != nil {
			switch err {
			case auth.ErrInvalidPassword:
				serveError(w, err.Error(), http.StatusUnauthorized)
			case auth.ErrUserNotConfirmed:
				serveError(w, err.Error(), http.StatusForbidden)
			case auth.ErrUserNotFound:
				serveError(w, err.Error(), http.StatusNotFound)
			case auth.ErrAccessTokenExpired:
				serveError(w, err.Error(), http.StatusGone)
			default:
				serveError(w, "internal error", http.StatusInternalServerError)
			}

			return
		}

		cookieAccessToken := &http.Cookie{
			Name:     "access_token",
			Value:    loginResp.AccessToken,
			HttpOnly: true,
			MaxAge:   int(svc.Cfg.AccessTokenTTL.Seconds()),
		}

		cookieRefreshToken := &http.Cookie{
			Name:     "refresh_token",
			Value:    loginResp.RefreshToken,
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
