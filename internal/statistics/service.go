package statistics

import (
	"context"

	"cafe-scheduling-api/pkg/apperror"
	"cafe-scheduling-api/pkg/timeutil"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Weekly(ctx context.Context, weekStart string) (WeeklyStatistics, error) {
	start, err := timeutil.ParseDate(weekStart)
	if err != nil {
		return WeeklyStatistics{}, apperror.New(apperror.ErrBadRequest, err.Error())
	}
	end := timeutil.EndOfWeek(start)
	items, err := s.repo.WeeklyHours(ctx, timeutil.FormatDate(start), timeutil.FormatDate(end))
	if err != nil {
		return WeeklyStatistics{}, err
	}
	return WeeklyStatistics{
		WeekStart: timeutil.FormatDate(start),
		WeekEnd:   timeutil.FormatDate(end),
		Employees: items,
	}, nil
}
