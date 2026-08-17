package roster

import (
	"context"
	"fmt"
	"sort"
	"time"
)

func BuildWeekView(ctx context.Context, weekStart string, assignments []Assignment) (WeekView, error) {
	if err := ctx.Err(); err != nil {
		return WeekView{}, err
	}
	start, err := time.Parse(dateLayout, weekStart)
	if err != nil {
		return WeekView{}, fmt.Errorf("validate week start: %w", invalidField("weekStart", "must use YYYY-MM-DD"))
	}
	if start.Weekday() != time.Monday {
		return WeekView{}, invalidField("weekStart", "must be a Monday")
	}

	view := WeekView{
		WeekStart: weekStart,
		WeekEnd:   formatDate(start.AddDate(0, 0, 6)),
		Days:      make([]Day, 7),
	}
	dayIndexes := make(map[string]int, 7)
	for i := range view.Days {
		date := formatDate(start.AddDate(0, 0, i))
		view.Days[i] = Day{Date: date, Assignments: make([]Assignment, 0)}
		dayIndexes[date] = i
	}

	ordered := make([]Assignment, len(assignments))
	for i, assignment := range assignments {
		if err := validateAssignment(ctx, assignment); err != nil {
			return WeekView{}, fmt.Errorf("assignment %d: %v", i, err)
		}
		ordered[i] = cloneAssignment(assignment)
	}
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].WorkDate != ordered[j].WorkDate {
			return ordered[i].WorkDate < ordered[j].WorkDate
		}
		if ordered[i].StartTime != ordered[j].StartTime {
			return ordered[i].StartTime < ordered[j].StartTime
		}
		return ordered[i].ID < ordered[j].ID
	})

	for _, assignment := range ordered {
		if err := ctx.Err(); err != nil {
			return WeekView{}, err
		}
		index, ok := dayIndexes[assignment.WorkDate]
		if !ok {
			return WeekView{}, invalidField("workDate", "must be inside the requested week")
		}
		hours := assignmentHours(assignment)
		view.Days[index].Assignments = append(view.Days[index].Assignments, cloneAssignment(assignment))
		view.Days[index].Hours += hours
		view.Hours += hours
	}
	return view, nil
}
