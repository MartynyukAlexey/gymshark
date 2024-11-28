package auth

import (
	"errors"
	"log/slog"

	"github.com/MartynyukAlexey/gymshark/internal/config"
	"github.com/MartynyukAlexey/gymshark/internal/smtp"
	"github.com/MartynyukAlexey/gymshark/internal/storage"
)

type Service struct {
	Storage *storage.Storage
	Mailer  *smtp.SMTPMailer
	Logger  *slog.Logger
	Cfg     *config.AuthConfig
}

var (
	// users
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidName     = errors.New("invalid name or surname")

	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserNotConfirmed     = errors.New("user is not confirmed yet")
	ErrUserAlreadyConfirmed = errors.New("user is already confirmed")

	// confirmation codes
	ErrInvalidCode = errors.New("invalid code")
	ErrCodeExpired = errors.New("code expired")

	// tokens
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrAccessTokenExpired  = errors.New("access token expired")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrRefreshTokenReuse   = errors.New("refresh token reuse")
)
