package businesshours

import (
	"context"
	"time"

	"cafe-scheduling-api/pkg/apperror"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]BusinessHours, error) {
	return s.repo.List(ctx)
}

func (s *Service) UpsertWeekly(ctx context.Context, items []UpsertItem) ([]BusinessHours, error) {
	if len(items) == 0 {
		return nil, apperror.New(apperror.ErrBadRequest, "at least one business hours item is required")
	}

	seen := make(map[int]struct{}, len(items))
	for _, item := range items {
		if item.DayOfWeek < 1 || item.DayOfWeek > 7 {
			return nil, apperror.New(apperror.ErrBadRequest, "dayOfWeek must be between 1 and 7")
		}
		if _, exists := seen[item.DayOfWeek]; exists {
			return nil, apperror.New(apperror.ErrBadRequest, "duplicate dayOfWeek in request")
		}
		seen[item.DayOfWeek] = struct{}{}

		if item.IsClosed {
			item.OpenTime = "00:00:00"
			item.CloseTime = "00:00:00"
		} else {
			open, close, err := normalizeTimes(item.OpenTime, item.CloseTime)
			if err != nil {
				return nil, err
			}
			item.OpenTime = open
			item.CloseTime = close
		}

		if err := s.repo.Upsert(ctx, BusinessHours{
			DayOfWeek: item.DayOfWeek,
			OpenTime:  item.OpenTime,
			CloseTime: item.CloseTime,
			IsClosed:  item.IsClosed,
		}); err != nil {
			return nil, err
		}
	}

	return s.repo.List(ctx)
}

func normalizeTimes(startRaw, endRaw string) (string, string, error) {
	start, err := time.Parse("15:04", startRaw)
	if err != nil {
		return "", "", apperror.New(apperror.ErrBadRequest, "openTime must be in HH:MM format")
	}
	end, err := time.Parse("15:04", endRaw)
	if err != nil {
		return "", "", apperror.New(apperror.ErrBadRequest, "closeTime must be in HH:MM format")
	}
	if !end.After(start) {
		return "", "", apperror.New(apperror.ErrBadRequest, "closeTime must be after openTime")
	}
	return start.Format("15:04:00"), end.Format("15:04:00"), nil
}
