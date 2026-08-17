package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrUserNotFound = errors.New("user not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (User, error) {
	const query = `
		SELECT id, username, password_hash, role, employee_id, created_at, updated_at
		FROM users
		WHERE username = ?`

	var u User
	var employeeID sql.NullInt64
	if err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.Role, &employeeID, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("query user by username: %w", err)
	}
	if employeeID.Valid {
		u.EmployeeID = &employeeID.Int64
	}
	return u, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (User, error) {
	const query = `
		SELECT id, username, password_hash, role, employee_id, created_at, updated_at
		FROM users
		WHERE id = ?`

	var u User
	var employeeID sql.NullInt64
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.Role, &employeeID, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("query user by id: %w", err)
	}
	if employeeID.Valid {
		u.EmployeeID = &employeeID.Int64
	}
	return u, nil
}

func (r *Repository) Create(ctx context.Context, u User) (int64, error) {
	const query = `
		INSERT INTO users (username, password_hash, role, employee_id)
		VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, u.Username, u.PasswordHash, u.Role, u.EmployeeID)
	if err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get user id: %w", err)
	}
	return id, nil
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return count, nil
}
