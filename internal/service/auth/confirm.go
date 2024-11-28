package auth

import (
	"net/mail"
	"time"

	"github.com/MartynyukAlexey/gymshark/internal/storage/models"
	"golang.org/x/crypto/bcrypt"
)

type ConfirmReq struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (s *Service) Confirm(req *ConfirmReq) error {
	if err := validateConfirmReq(req); err != nil {
		return err
	}

	user, err := s.Storage.User.GetByEmail(req.Email)
	if err != nil {
		if err == models.ErrUserNotFound {
			return ErrUserNotFound
		}

		s.Logger.Error("failed to get user by email", "err", err)
		return err
	}

	switch user.State {
	case models.UserStateActive:
		return ErrUserAlreadyConfirmed
	case models.UserStateDeleted:
		return ErrUserNotFound
	}

	codes, err := s.Storage.Code.GetAllByUser(user.ID, models.CodeScopeConfirm)
	if err != nil {
		s.Logger.Error("failed to get confirmation codes", "err", err)
		return err
	}

	for _, code := range codes {
		if err := bcrypt.CompareHashAndPassword(code.Hash, []byte(req.Code)); err != nil {
			if err != bcrypt.ErrMismatchedHashAndPassword {
				s.Logger.Error("failed to verify confirmation code", "err", err)
				return err
			}
		} else {
			// code matches

			if code.ExpiresAt.Before(time.Now()) {
				return ErrCodeExpired
			}

			if code.Scope != models.CodeScopeConfirm {
				return ErrInvalidPassword
			}

			if err := s.Storage.User.UpdateStatus(user.ID, models.UserStateActive); err != nil {
				s.Logger.Error("failed to activate user account", "err", err)
				return err
			}

			return nil
		}
	}

	return ErrInvalidCode
}

func validateConfirmReq(req *ConfirmReq) error {
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return ErrInvalidEmail
	}

	if len(req.Code) == 0 {
		return ErrInvalidCode
	}

	return nil
}
