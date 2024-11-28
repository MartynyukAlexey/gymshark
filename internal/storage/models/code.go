package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type CodeScope string

const (
	CodeScopeReset   CodeScope = "reset"
	CodeScopeConfirm CodeScope = "confirm"
)

type Code struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Hash   []byte

	Scope CodeScope

	CreatedAt time.Time
	ExpiresAt time.Time
}

var (
	ErrCodeNotFound = errors.New("code not found")
)
