package businesshours

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrBusinessHoursNotFound = errors.New("business hours not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]BusinessHours, error) {
	const query = `SELECT id, day_of_week, TIME_FORMAT(open_time, '%H:%i'), TIME_FORMAT(close_time, '%H:%i'), is_closed, created_at, updated_at
		FROM business_hours ORDER BY day_of_week`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query business hours: %w", err)
	}
	defer rows.Close()

	var items []BusinessHours
	for rows.Next() {
		var item BusinessHours
		if err := rows.Scan(&item.ID, &item.DayOfWeek, &item.OpenTime, &item.CloseTime, &item.IsClosed, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan business hours: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate business hours: %w", err)
	}
	return items, nil
}

func (r *Repository) GetByDay(ctx context.Context, day int) (BusinessHours, error) {
	const query = `SELECT id, day_of_week, TIME_FORMAT(open_time, '%H:%i'), TIME_FORMAT(close_time, '%H:%i'), is_closed, created_at, updated_at
		FROM business_hours WHERE day_of_week = ?`
	var item BusinessHours
	if err := r.db.QueryRowContext(ctx, query, day).Scan(
		&item.ID, &item.DayOfWeek, &item.OpenTime, &item.CloseTime, &item.IsClosed, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return BusinessHours{}, ErrBusinessHoursNotFound
		}
		return BusinessHours{}, fmt.Errorf("query business hours by day: %w", err)
	}
	return item, nil
}

func (r *Repository) Upsert(ctx context.Context, item BusinessHours) error {
	const query = `
		INSERT INTO business_hours (day_of_week, open_time, close_time, is_closed)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE open_time = VALUES(open_time), close_time = VALUES(close_time), is_closed = VALUES(is_closed)`
	_, err := r.db.ExecContext(ctx, query, item.DayOfWeek, item.OpenTime, item.CloseTime, item.IsClosed)
	if err != nil {
		return fmt.Errorf("upsert business hours: %w", err)
	}
	return nil
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM business_hours`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count business hours: %w", err)
	}
	return count, nil
}
