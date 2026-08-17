package statistics

type WeeklyStatistics struct {
	WeekStart string                `json:"weekStart"`
	WeekEnd   string                `json:"weekEnd"`
	Employees []EmployeeWeeklyHours `json:"employees"`
}

type EmployeeWeeklyHours struct {
	EmployeeID   int64   `json:"employeeId"`
	EmployeeName string  `json:"employeeName"`
	Hours        float64 `json:"hours"`
	ShiftCount   int     `json:"shiftCount"`
}
