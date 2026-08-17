package schedule

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrScheduleNotFound = errors.New("schedule not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const selectColumns = `
	s.id,
	s.employee_id,
	e.name,
	s.shift_type_id,
	st.name,
	DATE_FORMAT(s.work_date, '%Y-%m-%d'),
	TIME_FORMAT(s.start_time, '%H:%i'),
	TIME_FORMAT(s.end_time, '%H:%i'),
	s.created_at,
	s.updated_at`

func (r *Repository) ListByWeek(ctx context.Context, weekStart, weekEnd string, employeeID *int64) ([]Schedule, error) {
	query := `
		SELECT ` + selectColumns + `
		FROM schedules s
		JOIN employees e ON e.id = s.employee_id
		JOIN shift_types st ON st.id = s.shift_type_id
		WHERE s.work_date BETWEEN ? AND ?`
	args := []any{weekStart, weekEnd}
	if employeeID != nil {
		query += ` AND s.employee_id = ?`
		args = append(args, *employeeID)
	}
	query += ` ORDER BY s.work_date, s.start_time, s.id`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query schedules: %w", err)
	}
	defer rows.Close()

	var items []Schedule
	for rows.Next() {
		item, err := scanSchedule(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate schedules: %w", err)
	}
	return items, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (Schedule, error) {
	query := `
		SELECT ` + selectColumns + `
		FROM schedules s
		JOIN employees e ON e.id = s.employee_id
		JOIN shift_types st ON st.id = s.shift_type_id
		WHERE s.id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	item, err := scanSchedule(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Schedule{}, ErrScheduleNotFound
		}
		return Schedule{}, err
	}
	return item, nil
}

func (r *Repository) Create(ctx context.Context, s Schedule) (int64, error) {
	const query = `
		INSERT INTO schedules (employee_id, shift_type_id, work_date, start_time, end_time)
		VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, s.EmployeeID, s.ShiftTypeID, s.WorkDate, s.StartTime, s.EndTime)
	if err != nil {
		return 0, fmt.Errorf("insert schedule: %w", err)
	}
	return result.LastInsertId()
}

func (r *Repository) Update(ctx context.Context, id int64, s Schedule) error {
	const query = `
		UPDATE schedules
		SET employee_id = ?, shift_type_id = ?, work_date = ?, start_time = ?, end_time = ?
		WHERE id = ?`
	if _, err := r.db.ExecContext(ctx, query, s.EmployeeID, s.ShiftTypeID, s.WorkDate, s.StartTime, s.EndTime, id); err != nil {
		return fmt.Errorf("update schedule: %w", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM schedules WHERE id = ?`, id); err != nil {
		return fmt.Errorf("delete schedule: %w", err)
	}
	return nil
}

type overlapInfo struct {
	ID        int64
	StartTime string
	EndTime   string
}

func (r *Repository) FindOverlap(ctx context.Context, employeeID int64, workDate, startTime, endTime string, excludeID int64) (overlapInfo, error) {
	const query = `
		SELECT id, TIME_FORMAT(start_time, '%H:%i'), TIME_FORMAT(end_time, '%H:%i')
		FROM schedules
		WHERE employee_id = ?
		  AND work_date = ?
		  AND id <> ?
		  AND start_time < ?
		  AND end_time > ?
		LIMIT 1`
	var overlap overlapInfo
	err := r.db.QueryRowContext(ctx, query, employeeID, workDate, excludeID, endTime, startTime).Scan(
		&overlap.ID, &overlap.StartTime, &overlap.EndTime,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return overlapInfo{}, nil
	}
	if err != nil {
		return overlapInfo{}, fmt.Errorf("query overlapping schedule: %w", err)
	}
	return overlap, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanSchedule(row rowScanner) (Schedule, error) {
	var item Schedule
	if err := row.Scan(
		&item.ID,
		&item.EmployeeID,
		&item.EmployeeName,
		&item.ShiftTypeID,
		&item.ShiftTypeName,
		&item.WorkDate,
		&item.StartTime,
		&item.EndTime,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return Schedule{}, fmt.Errorf("scan schedule: %w", err)
	}
	item.Hours = durationHours(item.StartTime, item.EndTime)
	return item, nil
}

func durationHours(start, end string) float64 {
	startMinutes := timeToMinutes(start)
	endMinutes := timeToMinutes(end)
	if endMinutes <= startMinutes {
		return 0
	}
	return float64(endMinutes-startMinutes) / 60
}

func timeToMinutes(value string) int {
	var h, m int
	_, _ = fmt.Sscanf(value, "%d:%d", &h, &m)
	return h*60 + m
}
