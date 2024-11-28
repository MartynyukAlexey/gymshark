package storage

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"github.com/MartynyukAlexey/gymshark/internal/storage/models"
	"github.com/MartynyukAlexey/gymshark/internal/storage/postgres"
)

type Storage struct {
	User  UserStorage
	Code  CodeStorage
	Token TokenStorage
}

func NewStorage(db *sql.DB, _ *minio.Client) *Storage {
	return &Storage{
		User:  postgres.NewUserStorage(db),
		Code:  postgres.NewCodeStorage(db),
		Token: postgres.NewTokenStorage(db),
	}
}

type UserStorage interface {
	Insert(user *models.User) error

	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)

	UpdateStatus(id uuid.UUID, state models.UserState) error

	DeleteByEmail(email string) error
	DeleteByEmailIfInactive(email string) error
}

type CodeStorage interface {
	Insert(code *models.Code) error

	GetByID(id uuid.UUID) (*models.Code, error)
	GetAllByUser(userID uuid.UUID, scope models.CodeScope) ([]*models.Code, error)

	DeleteAllByUser(userID uuid.UUID) error
	DeleteAllExpired() error
}

type TokenStorage interface {
	Insert(token *models.Token) error
	// insert new token and mark parent as used
	InsertChild(parentID uuid.UUID, token *models.Token) error

	GetByID(id uuid.UUID) (*models.Token, error)
	GetAllByUser(userID uuid.UUID) ([]*models.Token, error)
	GetAllByUserScope(userID uuid.UUID, scope models.TokenScope) ([]*models.Token, error)

	DeleteAllByUser(userID uuid.UUID) error
	DeleteAllByBranch(userID uuid.UUID, branch uuid.UUID) error
}
