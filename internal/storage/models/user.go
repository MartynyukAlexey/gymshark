package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserState string

const (
	UserStatePending UserState = "pending"
	UserStateActive  UserState = "active"
	UserStateDeleted UserState = "deleted"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash []byte
	State        UserState

	AvatarID  string
	FirstName string
	LastName  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrDuplicateEmail = errors.New("email is already taken")
)
