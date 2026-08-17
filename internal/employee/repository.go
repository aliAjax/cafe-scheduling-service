package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrEmployeeNotFound = errors.New("employee not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]Employee, error) {
	const query = `
		SELECT e.id, e.name, e.phone, e.position, e.active, u.id, e.created_at, e.updated_at
		FROM employees e
		LEFT JOIN users u ON u.employee_id = e.id
		ORDER BY e.id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query employees: %w", err)
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var e Employee
		var userID sql.NullInt64
		if err := rows.Scan(&e.ID, &e.Name, &e.Phone, &e.Position, &e.Active, &userID, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan employee: %w", err)
		}
		if userID.Valid {
			e.UserID = &userID.Int64
		}
		employees = append(employees, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate employees: %w", err)
	}
	return employees, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (Employee, error) {
	const query = `
		SELECT e.id, e.name, e.phone, e.position, e.active, u.id, e.created_at, e.updated_at
		FROM employees e
		LEFT JOIN users u ON u.employee_id = e.id
		WHERE e.id = ?`

	var e Employee
	var userID sql.NullInt64
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.Name, &e.Phone, &e.Position, &e.Active, &userID, &e.CreatedAt, &e.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Employee{}, ErrEmployeeNotFound
		}
		return Employee{}, fmt.Errorf("query employee by id: %w", err)
	}
	if userID.Valid {
		e.UserID = &userID.Int64
	}
	return e, nil
}

func (r *Repository) Create(ctx context.Context, e Employee) (int64, error) {
	const query = `INSERT INTO employees (name, phone, position) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, e.Name, e.Phone, e.Position)
	if err != nil {
		return 0, fmt.Errorf("insert employee: %w", err)
	}
	return result.LastInsertId()
}

func (r *Repository) Update(ctx context.Context, id int64, fields map[string]any) error {
	query := `UPDATE employees SET `
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
	if first {
		return nil
	}
	query += ` WHERE id = ?`
	args = append(args, id)

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("update employee: %w", err)
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE employees SET active = 0 WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete employee: %w", err)
	}
	return nil
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM employees`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count employees: %w", err)
	}
	return count, nil
}
