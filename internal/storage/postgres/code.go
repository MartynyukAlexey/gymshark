package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MartynyukAlexey/gymshark/internal/storage/models"
	"github.com/google/uuid"
)

type CodeStorage struct {
	db *sql.DB
}

func NewCodeStorage(db *sql.DB) *CodeStorage {
	return &CodeStorage{
		db: db,
	}
}

func (s *CodeStorage) Insert(code *models.Code) error {
	stmt := `
		INSERT INTO codes (
			user_id, hash, scope, expires_at
		) VALUES (
			$1, $2, $3, $4
		) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		stmt,
		code.UserID,
		code.Hash,
		code.Scope,
		code.ExpiresAt,
	).Scan(
		&code.ID,
		&code.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert code: %w", err)
	}

	return nil
}

func (s *CodeStorage) GetByID(id uuid.UUID) (*models.Code, error) {
	return nil, nil
}

func (s *CodeStorage) GetAllByUser(userID uuid.UUID, scope models.CodeScope) ([]*models.Code, error) {
	stmt := `
		SELECT
			id,
			user_id,
			hash,
			scope,
			expires_at,
			created_at
		FROM codes
		WHERE user_id = $1 AND scope = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, stmt, userID, scope)
	if err != nil {
		return nil, fmt.Errorf("failed to get codes for user: %w", err)
	}
	defer rows.Close()

	var codes []*models.Code
	for rows.Next() {
		var code models.Code
		if err := rows.Scan(
			&code.ID,
			&code.UserID,
			&code.Hash,
			&code.Scope,
			&code.ExpiresAt,
			&code.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to get codes for user: %w", err)
		}

		codes = append(codes, &code)
	}

	return codes, nil
}

func (s *CodeStorage) DeleteAllByUser(userID uuid.UUID) error {
	stmt := `
		DELETE FROM codes
		WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, stmt, userID)
	if err != nil {
		return fmt.Errorf("failed to delete codes for user: %w", err)
	}

	return nil
}

func (s *CodeStorage) DeleteAllExpired() error {
	stmt := `
		DELETE FROM codes
		WHERE expires_at < NOW()
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, stmt)
	if err != nil {
		return fmt.Errorf("failed to delete expired codes: %w", err)
	}

	return nil
}
