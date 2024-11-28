package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MartynyukAlexey/gymshark/internal/storage/models"
	"github.com/google/uuid"
)

type TokenStorage struct {
	db *sql.DB
}

func NewTokenStorage(db *sql.DB) *TokenStorage {
	return &TokenStorage{
		db: db,
	}
}

func (s *TokenStorage) Insert(token *models.Token) error {
	return nil
}

func (s *TokenStorage) InsertChild(parentID uuid.UUID, token *models.Token) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	stmt := `UPDATE tokens SET status = $1 WHERE id = $2 RETURNING branch`

	row := tx.QueryRowContext(ctx, stmt, models.TokenStatusUsed, parentID)

	var branch uuid.UUID
	err = row.Scan(&branch)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ErrTokenNotFound
		}

		return fmt.Errorf("failed to get branch by parent token id: %w", err)
	}

	stmt = `INSERT INTO tokens (user_id, hash, branch, status, scope, created_at, expires_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = tx.ExecContext(ctx, stmt,
		token.UserID,
		token.Hash,
		token.Branch,
		token.Status,
		token.Scope,
		token.CreatedAt,
		token.ExpiresAt,
	)

	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert child token in transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *TokenStorage) GetByID(id uuid.UUID) (*models.Token, error) {
	return nil, nil
}

func (s *TokenStorage) GetAllByUser(userID uuid.UUID) ([]*models.Token, error) {
	return nil, nil
}

func (s *TokenStorage) GetAllByUserScope(userID uuid.UUID, scope models.TokenScope) ([]*models.Token, error) {
	return nil, nil
}

func (s *TokenStorage) DeleteAllByUser(userID uuid.UUID) error {
	return nil
}

func (s *TokenStorage) DeleteAllByBranch(userID uuid.UUID, branch uuid.UUID) error {
	return nil
}
