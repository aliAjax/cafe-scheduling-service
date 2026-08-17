package statistics

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) WeeklyHours(ctx context.Context, weekStart, weekEnd string) ([]EmployeeWeeklyHours, error) {
	const query = `
		SELECT
			e.id,
			e.name,
			COALESCE(SUM(TIME_TO_SEC(TIMEDIFF(s.end_time, s.start_time))) / 3600.0, 0),
			COUNT(s.id)
		FROM employees e
		LEFT JOIN schedules s
			ON s.employee_id = e.id
		   AND s.work_date BETWEEN ? AND ?
		WHERE e.active = 1
		GROUP BY e.id, e.name
		ORDER BY e.id`

	rows, err := r.db.QueryContext(ctx, query, weekStart, weekEnd)
	if err != nil {
		return nil, fmt.Errorf("query weekly statistics: %w", err)
	}
	defer rows.Close()

	var items []EmployeeWeeklyHours
	for rows.Next() {
		var item EmployeeWeeklyHours
		var hours sql.NullFloat64
		if err := rows.Scan(&item.EmployeeID, &item.EmployeeName, &hours, &item.ShiftCount); err != nil {
			return nil, fmt.Errorf("scan weekly statistics: %w", err)
		}
		if hours.Valid {
			item.Hours = hours.Float64
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate weekly statistics: %w", err)
	}
	return items, nil
}
