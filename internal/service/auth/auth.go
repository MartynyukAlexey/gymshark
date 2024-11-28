package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func (s *Service) Authorize(accessToken string) (uuid.UUID, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidAccessToken
		}
		return s.Cfg.JWTKey, nil
	})

	if err != nil {
		s.Logger.Error("failed to parse access token", "err", err)
		return uuid.Nil, ErrInvalidAccessToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, ErrInvalidAccessToken
	}

	expirationClaim, ok := claims["exp"]
	if !ok {
		return uuid.Nil, ErrInvalidAccessToken
	}

	expirationFloat, ok := expirationClaim.(float64)
	if !ok {
		return uuid.Nil, ErrInvalidAccessToken
	}

	if int64(expirationFloat) < time.Now().Unix() {
		return uuid.Nil, ErrAccessTokenExpired
	}

	userIDClaim, ok := claims["sub"]
	if !ok {
		return uuid.Nil, ErrInvalidAccessToken
	}

	userIDString, ok := userIDClaim.(string)
	if !ok {
		return uuid.Nil, ErrInvalidAccessToken
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, ErrInvalidAccessToken
	}

	return userID, nil
}
