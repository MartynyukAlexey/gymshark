package auth

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/MartynyukAlexey/gymshark/internal/storage/models"
)

type RegisterReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (s *Service) Register(req *RegisterReq) (uuid.UUID, error) {
	if err := validateRegisterReq(req); err != nil {
		return uuid.Nil, err
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("failed to hash password", "err", err)
		return uuid.Nil, err
	}

	m := &models.User{
		Email:        req.Email,
		PasswordHash: passHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	// remove deleted accounts and accounts from unsuccessfult registrations
	if err := s.Storage.User.DeleteByEmailIfInactive(req.Email); err != nil {
		s.Logger.Error("failed to delete old non-active user", "err", err)
		return uuid.Nil, err
	}

	if err := s.Storage.User.Insert(m); err != nil {
		if err == models.ErrDuplicateEmail {
			return uuid.Nil, ErrUserAlreadyExists
		}

		s.Logger.Error("failed to save user", "err", err)
		return uuid.Nil, err
	}

	code, err := generateCode(8)
	if err != nil {
		s.Logger.Error("failed to generate activation code", "err", err)
		return uuid.Nil, err
	}

	codeHash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("failed to hash activation code", "err", err)
		return uuid.Nil, err
	}

	t := &models.Code{
		UserID:    m.ID,
		Hash:      codeHash,
		Scope:     models.CodeScopeConfirm,
		ExpiresAt: m.CreatedAt.Add(24 * time.Hour),
	}

	if err := s.Storage.Code.Insert(t); err != nil {
		s.Logger.Error("failed to save activation code", "err", err)
		return uuid.Nil, err
	}

	go func() {
		if err := s.Mailer.SendActivationEmail(m.Email, code); err != nil {
			s.Logger.Error("failed to send activation email", "err", err)
		}
	}()

	return m.ID, nil
}

func validateRegisterReq(req *RegisterReq) error {
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return ErrInvalidEmail
	}

	if len(req.Password) < 8 {
		return ErrInvalidPassword
	}

	if len(req.FirstName) == 0 || len(req.LastName) == 0 {
		return ErrInvalidName
	}

	return nil
}
