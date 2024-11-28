package service

import (
	"log/slog"

	"github.com/MartynyukAlexey/gymshark/internal/config"
	"github.com/MartynyukAlexey/gymshark/internal/service/auth"
	"github.com/MartynyukAlexey/gymshark/internal/service/user"
	"github.com/MartynyukAlexey/gymshark/internal/smtp"
	"github.com/MartynyukAlexey/gymshark/internal/storage"
)

type Service struct {
	User *user.Service
	Auth *auth.Service
}

type ServiceOpts struct {
	Storage    *storage.Storage
	Mailer     *smtp.SMTPMailer
	Logger     *slog.Logger
	AuthConfig *config.AuthConfig
}

func NewService(opts *ServiceOpts) *Service {
	return &Service{
		Auth: &auth.Service{
			Storage: opts.Storage,
			Mailer:  opts.Mailer,
			Logger:  opts.Logger,
			Cfg:     opts.AuthConfig,
		},

		User: &user.Service{},
	}
}
