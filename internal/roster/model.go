package roster

import "time"

type Assignment struct {
	ID         int64    `json:"id"`
	EmployeeID int64    `json:"employeeId"`
	WorkDate   string   `json:"workDate"`
	StartTime  string   `json:"startTime"`
	EndTime    string   `json:"endTime"`
	Note       *string  `json:"note,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

type Day struct {
	Date        string       `json:"date"`
	Assignments []Assignment `json:"assignments"`
	Hours       float64      `json:"hours"`
}

type WeekView struct {
	WeekStart string  `json:"weekStart"`
	WeekEnd   string  `json:"weekEnd"`
	Days      []Day   `json:"days"`
	Hours     float64 `json:"hours"`
}

func cloneAssignment(a Assignment) Assignment {
	copyOf := a
	if a.Note != nil {
		note := *a.Note
		copyOf.Note = &note
	}
	copyOf.Tags = append([]string(nil), a.Tags...)
	return copyOf
}

func cloneDay(day Day) Day {
	copyOf := day
	copyOf.Assignments = make([]Assignment, len(day.Assignments))
	for i, assignment := range day.Assignments {
		copyOf.Assignments[i] = cloneAssignment(assignment)
	}
	return copyOf
}

func cloneWeek(view WeekView) WeekView {
	copyOf := view
	copyOf.Days = make([]Day, len(view.Days))
	for i, day := range view.Days {
		copyOf.Days[i] = cloneDay(day)
	}
	return copyOf
}

func formatDate(value time.Time) string {
	return value.Format("2006-01-02")
}
