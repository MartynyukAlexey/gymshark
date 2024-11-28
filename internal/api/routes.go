package api

import (
	"log/slog"
	"net/http"

	"github.com/MartynyukAlexey/gymshark/internal/api/auth"
	"github.com/MartynyukAlexey/gymshark/internal/service"
)

func Routes(service *service.Service, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	m := middlewareEnv{
		svc:    service.Auth,
		logger: logger,
	}

	mux.Handle("POST /api/v1/register", auth.HandleRegistration(service.Auth, logger))
	mux.Handle("POST /api/v1/confirm", auth.HandleConfirmation(service.Auth, logger))
	mux.Handle("POST /api/v1/login", auth.HandleLogin(service.Auth, logger))
	mux.Handle("POST /api/v1/refresh", auth.HandleRefresh(service.Auth, logger))

	mux.Handle("GET /api/v1/test", m.RequireAuth(auth.HandleTest(service.Auth, logger)))

	return mux
}
