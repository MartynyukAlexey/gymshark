package auth

import (
	"time"

	"github.com/MartynyukAlexey/gymshark/internal/storage/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RefreshReq struct {
	Email        string `json:"email"`
	RefreshToken string `json:"access_token"`
}

type RefreshResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) Refresh(req *RefreshReq) (RefreshResp, error) {
	if err := validateRefreshReq(req); err != nil {
		return RefreshResp{}, err
	}

	user, err := s.Storage.User.GetByEmail(req.Email)
	if err != nil {
		if err == models.ErrUserNotFound {
			return RefreshResp{}, ErrUserNotFound
		}

		s.Logger.Error("failed to get user by email", "err", err)
		return RefreshResp{}, err
	}

	tokens, err := s.Storage.Token.GetAllByUserScope(user.ID, models.TokenScopeRefresh)
	if err != nil {
		s.Logger.Error("failed to get tokens for user", "err", err)
		return RefreshResp{}, err
	}

	for _, token := range tokens {
		if err := bcrypt.CompareHashAndPassword(token.Hash, []byte(req.RefreshToken)); err != nil {
			if err != bcrypt.ErrMismatchedHashAndPassword {
				s.Logger.Error("failed to verify refrest token", "err", err)
				return RefreshResp{}, err
			}
		} else {
			if token.ExpiresAt.Before(time.Now()) {
				return RefreshResp{}, ErrRefreshTokenExpired
			}

			// an attempt to reuse the token (suspect token leakage).
			// revoking all tokens in the branch
			if token.Status != models.TokenStatusActive {
				if err := s.Storage.Token.DeleteAllByBranch(user.ID, token.Branch); err != nil {
					s.Logger.Error("failed to delete tokens", "err", err)
				}

				return RefreshResp{}, ErrInvalidRefreshToken
			}

			// successful refresh
			newRefreshToken := uuid.New().String()
			if err := s.Storage.Token.InsertChild(token.ID, &models.Token{
				UserID:    user.ID,
				Hash:      []byte(newRefreshToken),
				Branch:    token.Branch,
				Status:    models.TokenStatusActive,
				Scope:     models.TokenScopeRefresh,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(time.Minute * 30),
			}); err != nil {
				s.Logger.Error("failed to insert new child token", "err", err)
				return RefreshResp{}, err
			}

			newAccessToken, err := generateAccessToken(user.ID, s.Cfg.AccessTokenTTL, s.Cfg.JWTKey)
			if err != nil {
				s.Logger.Error("failed to generate access token", "err", err)
				return RefreshResp{}, err
			}

			return RefreshResp{
				AccessToken:  newAccessToken,
				RefreshToken: newRefreshToken,
			}, nil
		}
	}

	return RefreshResp{}, ErrInvalidRefreshToken
}

func validateRefreshReq(req *RefreshReq) error {
	return nil
}
