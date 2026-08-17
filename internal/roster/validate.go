package roster

import (
	"context"
	"fmt"
	"time"
)

const dateLayout = "2006-01-02"

func validateAssignment(ctx context.Context, assignment Assignment) error {
	if assignment.ID <= 0 {
		return invalidField("id", "must be positive")
	}
	if assignment.EmployeeID <= 0 {
		return invalidField("employeeId", "must be positive")
	}
	if _, err := time.Parse(dateLayout, assignment.WorkDate); err != nil {
		return invalidField("workDate", "must use YYYY-MM-DD")
	}
	start, err := parseClock(assignment.StartTime)
	if err != nil {
		return fmt.Errorf("validate start time: %w", invalidField("startTime", err.Error()))
	}
	end, err := parseClock(assignment.EndTime)
	if err != nil {
		return fmt.Errorf("validate end time: %w", invalidField("endTime", err.Error()))
	}
	if start == end {
		return invalidField("endTime", "must differ from startTime")
	}
	return nil
}

func parseClock(raw string) (int, error) {
	parsed, err := time.Parse("15:04", raw)
	if err != nil {
		return 0, fmt.Errorf("must use HH:MM")
	}
	return parsed.Hour()*60 + parsed.Minute(), nil
}

func assignmentHours(assignment Assignment) float64 {
	start, _ := parseClock(assignment.StartTime)
	end, _ := parseClock(assignment.EndTime)
	if end < start {
		end += 24 * 60
	}
	return float64(end-start) / 60
}
