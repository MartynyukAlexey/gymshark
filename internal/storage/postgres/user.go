package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/MartynyukAlexey/gymshark/internal/storage/models"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) Insert(user *models.User) error {
	stmt := `
        INSERT INTO "users" (
            email, password, avatar_id, first_name, last_name
        ) VALUES (
            $1, $2, $3, $4, $5
        ) RETURNING id, created_at, updated_at
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, stmt,
		user.Email,
		user.PasswordHash,
		user.AvatarID,
		user.FirstName,
		user.LastName,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" && err.Constraint == "users_email_key" {
				return models.ErrDuplicateEmail
			}
		}

		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (s *UserStorage) GetByID(id uuid.UUID) (*models.User, error) {
	stmt := `
	    SELECT
			id,
			email,
			password,
			avatar_id,
			first_name,
			last_name,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	err := s.db.QueryRowContext(ctx, stmt, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.AvatarID,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (s *UserStorage) GetByEmail(email string) (*models.User, error) {
	stmt := `
	    SELECT
			id,
			email,
			password,
			avatar_id,
			first_name,
			last_name,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	err := s.db.QueryRowContext(ctx, stmt, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.AvatarID,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (s *UserStorage) UpdateStatus(id uuid.UUID, state models.UserState) error {
	stmt := `
		UPDATE users
		SET state = $2
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := s.db.ExecContext(ctx, stmt, id, state)
	if err != nil {
		return fmt.Errorf("failed to execute update user status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to execute update user status: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrUserNotFound
	}

	return nil
}

func (s *UserStorage) DeleteByEmail(email string) error {
	stmt := `
		DELETE FROM users
		WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, stmt, email)
	if err != nil {
		return fmt.Errorf("failed to delete user by email: %w", err)
	}

	return nil
}

func (s *UserStorage) DeleteByEmailIfInactive(email string) error {
	stmt := `
		DELETE FROM users
		WHERE email = $1 AND state != $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, stmt, email, models.UserStateActive)
	if err != nil {
		return fmt.Errorf("failed to delete user by email: %w", err)
	}

	return nil
}
