package timeutil

import (
	"fmt"
	"time"
)

const DateLayout = "2006-01-02"

func ParseDate(value string) (time.Time, error) {
	parsed, err := time.Parse(DateLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("date must be in YYYY-MM-DD format")
	}
	return parsed, nil
}

func FormatDate(t time.Time) string {
	return t.Format(DateLayout)
}

func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	offset := (weekday + 6) % 7
	return time.Date(t.Year(), t.Month(), t.Day()-offset, 0, 0, 0, 0, t.Location())
}

func EndOfWeek(weekStart time.Time) time.Time {
	return weekStart.AddDate(0, 0, 6)
}

func DayOfWeekMondayOne(t time.Time) int {
	return (int(t.Weekday())+6)%7 + 1
}
