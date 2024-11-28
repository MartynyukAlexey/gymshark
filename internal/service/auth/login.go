package auth

import (
	"net/mail"
	"time"

	"github.com/MartynyukAlexey/gymshark/internal/storage/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) Login(req *LoginReq) (LoginResp, error) {
	if err := validateLoginReq(req); err != nil {
		return LoginResp{}, err
	}

	user, err := s.Storage.User.GetByEmail(req.Email)
	if err != nil {
		if err == models.ErrUserNotFound {
			return LoginResp{}, ErrUserNotFound
		}

		s.Logger.Error("/service/auth/login: failed to get user by email", "err", err)
		return LoginResp{}, err
	}

	switch user.State {
	case models.UserStatePending:
		return LoginResp{}, ErrUserNotConfirmed
	case models.UserStateDeleted:
		return LoginResp{}, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(req.Password)); err != nil {
		if err != bcrypt.ErrMismatchedHashAndPassword {
			s.Logger.Error("/service/auth/login: failed to verify password", "err", err)
			return LoginResp{}, err
		}

		return LoginResp{}, ErrPasswordMismatch
	}

	accessToken, err := generateAccessToken(user.ID, s.Cfg.AccessTokenTTL, s.Cfg.JWTKey)
	if err != nil {
		s.Logger.Error("/service/auth/login: failed to generate jwt token", "err", err)
		return LoginResp{}, err
	}

	// TODO: encrypting
	refreshToken := uuid.New().String()
	if err = s.Storage.Token.Insert(&models.Token{
		UserID:    user.ID,
		Hash:      []byte(refreshToken),
		Branch:    uuid.New(),
		Scope:     models.TokenScopeRefresh,
		ExpiresAt: time.Now().Add(s.Cfg.RefreshTokenTTL),
	}); err != nil {
		s.Logger.Error("/service/auth/login: failed to save refresh token", "user id", user.ID, "err", err)
		return LoginResp{}, err
	}

	return LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func validateLoginReq(req *LoginReq) error {
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return ErrInvalidEmail
	}

	if len(req.Password) < 8 {
		return ErrInvalidPassword
	}

	return nil
}
