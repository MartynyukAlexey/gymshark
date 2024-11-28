package auth

import (
	"crypto/rand"
	"encoding/base32"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func generateCode(bytesUsed int) (string, error) {
	b := make([]byte, bytesUsed)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b), nil
}

func generateAccessToken(userID uuid.UUID, ttl time.Duration, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "gymshark",
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString(key)
}
