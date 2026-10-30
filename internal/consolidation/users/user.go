package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// User represents a consolidation user
type User struct {
	ID         int
	Username   string
	CreatedAt  time.Time
	LastActive time.Time
}

// Repository interface for user operations
type Repository interface {
	// CreateUser creates a new user and returns the user ID
	CreateUser(ctx context.Context, username string) (*User, error)
	
	// GetUser retrieves a user by ID
	GetUser(ctx context.Context, id int) (*User, error)
	
	// GetUserByUsername retrieves a user by username
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	
	// ListUsers returns all users in the database
	ListUsers(ctx context.Context) ([]User, error)
	
	// UpdateLastActive updates the last_active timestamp for a user
	UpdateLastActive(ctx context.Context, userID int) error
	
	// UserExists checks if a user with the given username exists
	UserExists(ctx context.Context, username string) (bool, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL-backed user repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateUser creates a new user
func (r *PostgresRepository) CreateUser(ctx context.Context, username string) (*User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	user := &User{
		Username:   username,
		CreatedAt:  time.Now(),
		LastActive: time.Now(),
	}

	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO consolidation_users (username, created_at, last_active)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at, last_active`,
		username, user.CreatedAt, user.LastActive,
	).Scan(&user.ID, &user.CreatedAt, &user.LastActive)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (r *PostgresRepository) GetUser(ctx context.Context, id int) (*User, error) {
	user := &User{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, created_at, last_active
		 FROM consolidation_users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.LastActive)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: id=%d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, created_at, last_active
		 FROM consolidation_users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.LastActive)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: username=%s", username)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// ListUsers returns all users
func (r *PostgresRepository) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, username, created_at, last_active
		 FROM consolidation_users
		 ORDER BY last_active DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt, &user.LastActive); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// UpdateLastActive updates the last_active timestamp
func (r *PostgresRepository) UpdateLastActive(ctx context.Context, userID int) error {
	now := time.Now()
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE consolidation_users SET last_active = $1 WHERE id = $2`,
		now, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update last_active: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("user not found: id=%d", userID)
	}

	return nil
}

// UserExists checks if a user exists
func (r *PostgresRepository) UserExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM consolidation_users WHERE username = $1)`,
		username,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}
