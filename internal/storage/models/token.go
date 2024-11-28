package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TokenStatus string

const (
	TokenStatusActive  TokenStatus = "active"
	TokenStatusRevoked TokenStatus = "revoked"
	TokenStatusUsed    TokenStatus = "used"
)

type TokenScope string

const (
	TokenScopeAccess  TokenScope = "access"
	TokenScopeRefresh TokenScope = "refresh"
)

type Token struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Hash   []byte

	Branch uuid.UUID
	Status TokenStatus
	Scope  TokenScope

	CreatedAt time.Time
	ExpiresAt time.Time
}

var (
	ErrTokenNotFound = errors.New("token not found")
)
