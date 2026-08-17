package shift

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrShiftTypeNotFound = errors.New("shift type not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]ShiftType, error) {
	const query = `SELECT id, name, TIME_FORMAT(start_time, '%H:%i'), TIME_FORMAT(end_time, '%H:%i'), color, created_at, updated_at
		FROM shift_types ORDER BY start_time, id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query shift types: %w", err)
	}
	defer rows.Close()

	var items []ShiftType
	for rows.Next() {
		var s ShiftType
		if err := rows.Scan(&s.ID, &s.Name, &s.StartTime, &s.EndTime, &s.Color, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan shift type: %w", err)
		}
		items = append(items, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate shift types: %w", err)
	}
	return items, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (ShiftType, error) {
	const query = `SELECT id, name, TIME_FORMAT(start_time, '%H:%i'), TIME_FORMAT(end_time, '%H:%i'), color, created_at, updated_at
		FROM shift_types WHERE id = ?`
	var s ShiftType
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.Name, &s.StartTime, &s.EndTime, &s.Color, &s.CreatedAt, &s.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ShiftType{}, ErrShiftTypeNotFound
		}
		return ShiftType{}, fmt.Errorf("query shift type by id: %w", err)
	}
	return s, nil
}

func (r *Repository) Create(ctx context.Context, s ShiftType) (int64, error) {
	const query = `INSERT INTO shift_types (name, start_time, end_time, color) VALUES (?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, s.Name, s.StartTime, s.EndTime, s.Color)
	if err != nil {
		return 0, fmt.Errorf("insert shift type: %w", err)
	}
	return result.LastInsertId()
}

func (r *Repository) Update(ctx context.Context, id int64, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}
	query := `UPDATE shift_types SET `
	args := make([]any, 0, len(fields)+1)
	first := true
	for column, value := range fields {
		if !first {
			query += ", "
		}
		query += column + " = ?"
		args = append(args, value)
		first = false
	}
	query += ` WHERE id = ?`
	args = append(args, id)
	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("update shift type: %w", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM shift_types WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete shift type: %w", err)
	}
	return nil
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM shift_types`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count shift types: %w", err)
	}
	return count, nil
}
