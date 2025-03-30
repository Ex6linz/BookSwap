package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Ex6linz/BookSwap/backend/internal/auth/domain"
	_ "github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users 
		(id, name, email, password_hash, location, bio, avatar_url, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Location,
		user.Bio,
		user.AvatarURL,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrEmailExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, location, bio, avatar_url, rating, 
		created_at, updated_at FROM users WHERE email = $1`

	var user domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Location,
		&user.Bio,
		&user.AvatarURL,
		&user.Rating,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func isDuplicateKeyError(err error) bool {
	const uniqueViolationCode = "23505"
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == uniqueViolationCode
	}
	return false
}
