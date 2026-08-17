package schedule

import "time"

type Schedule struct {
	ID            int64     `json:"id"`
	EmployeeID    int64     `json:"employeeId"`
	EmployeeName  string    `json:"employeeName"`
	ShiftTypeID   int64     `json:"shiftTypeId"`
	ShiftTypeName string    `json:"shiftTypeName"`
	WorkDate      string    `json:"workDate"`
	StartTime     string    `json:"startTime"`
	EndTime       string    `json:"endTime"`
	Hours         float64   `json:"hours"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type CreateRequest struct {
	EmployeeID  int64  `json:"employeeId"`
	ShiftTypeID int64  `json:"shiftTypeId"`
	WorkDate    string `json:"workDate"`
}
