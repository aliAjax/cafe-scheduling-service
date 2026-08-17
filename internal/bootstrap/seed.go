package bootstrap

import (
	"context"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Seed(ctx context.Context, db *sql.DB) error {
	if err := seedEmployees(ctx, db); err != nil {
		return err
	}
	if err := seedUsers(ctx, db); err != nil {
		return err
	}
	if err := seedShiftTypes(ctx, db); err != nil {
		return err
	}
	if err := seedBusinessHours(ctx, db); err != nil {
		return err
	}
	return nil
}

func seedEmployees(ctx context.Context, db *sql.DB) error {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM employees`).Scan(&count); err != nil {
		return fmt.Errorf("count employees for seed: %w", err)
	}
	if count > 0 {
		return nil
	}

	employees := []struct {
		name     string
		phone    string
		position string
	}{
		{name: "Alice Chen", phone: "13800000001", position: "Barista"},
		{name: "Bob Li", phone: "13800000002", position: "Cashier"},
	}
	for _, item := range employees {
		if _, err := db.ExecContext(ctx,
			`INSERT INTO employees (name, phone, position) VALUES (?, ?, ?)`,
			item.name, item.phone, item.position,
		); err != nil {
			return fmt.Errorf("seed employee %s: %w", item.name, err)
		}
	}
	return nil
}

func seedUsers(ctx context.Context, db *sql.DB) error {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return fmt.Errorf("count users for seed: %w", err)
	}
	if count > 0 {
		return nil
	}

	managerHash, err := hashPassword("manager123")
	if err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx,
		`INSERT INTO users (username, password_hash, role, employee_id) VALUES (?, ?, 'manager', NULL)`,
		"manager", managerHash,
	); err != nil {
		return fmt.Errorf("seed manager user: %w", err)
	}

	rows, err := db.QueryContext(ctx, `SELECT id, name FROM employees ORDER BY id LIMIT 2`)
	if err != nil {
		return fmt.Errorf("query employees for user seed: %w", err)
	}
	type employeeRow struct {
		id   int64
		name string
	}
	var employees []employeeRow
	for rows.Next() {
		var item employeeRow
		if err := rows.Scan(&item.id, &item.name); err != nil {
			rows.Close()
			return fmt.Errorf("scan employee for user seed: %w", err)
		}
		employees = append(employees, item)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate employees for user seed: %w", err)
	}

	employeeHash, err := hashPassword("employee123")
	if err != nil {
		return err
	}
	for _, item := range employees {
		username := "alice"
		if item.id != employees[0].id {
			username = "bob"
		}
		if _, err := db.ExecContext(ctx,
			`INSERT INTO users (username, password_hash, role, employee_id) VALUES (?, ?, 'employee', ?)`,
			username, employeeHash, item.id,
		); err != nil {
			return fmt.Errorf("seed employee user %s: %w", username, err)
		}
	}
	return nil
}

func seedShiftTypes(ctx context.Context, db *sql.DB) error {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM shift_types`).Scan(&count); err != nil {
		return fmt.Errorf("count shift types for seed: %w", err)
	}
	if count > 0 {
		return nil
	}

	shiftTypes := []struct {
		name      string
		startTime string
		endTime   string
		color     string
	}{
		{name: "早班", startTime: "08:00:00", endTime: "14:00:00", color: "#4f8a8b"},
		{name: "晚班", startTime: "14:00:00", endTime: "22:00:00", color: "#d18f52"},
	}
	for _, item := range shiftTypes {
		if _, err := db.ExecContext(ctx,
			`INSERT INTO shift_types (name, start_time, end_time, color) VALUES (?, ?, ?, ?)`,
			item.name, item.startTime, item.endTime, item.color,
		); err != nil {
			return fmt.Errorf("seed shift type %s: %w", item.name, err)
		}
	}
	return nil
}

func seedBusinessHours(ctx context.Context, db *sql.DB) error {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM business_hours`).Scan(&count); err != nil {
		return fmt.Errorf("count business hours for seed: %w", err)
	}
	if count > 0 {
		return nil
	}

	for day := 1; day <= 7; day++ {
		if _, err := db.ExecContext(ctx,
			`INSERT INTO business_hours (day_of_week, open_time, close_time, is_closed) VALUES (?, '08:00:00', '22:00:00', 0)`,
			day,
		); err != nil {
			return fmt.Errorf("seed business hours for day %d: %w", day, err)
		}
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash seed password: %w", err)
	}
	return string(hash), nil
}
