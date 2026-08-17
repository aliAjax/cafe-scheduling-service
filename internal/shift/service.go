package shift

import (
	"context"
	"errors"
	"strings"
	"time"

	"cafe-scheduling-api/pkg/apperror"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]ShiftType, error) {
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id int64) (ShiftType, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrShiftTypeNotFound) {
			return ShiftType{}, apperror.New(apperror.ErrNotFound, "shift type not found")
		}
		return ShiftType{}, err
	}
	return item, nil
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (ShiftType, error) {
	start, end, err := normalizeTimes(req.StartTime, req.EndTime)
	if err != nil {
		return ShiftType{}, err
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return ShiftType{}, apperror.New(apperror.ErrBadRequest, "name is required")
	}

	id, err := s.repo.Create(ctx, ShiftType{
		Name:      req.Name,
		StartTime: start,
		EndTime:   end,
		Color:     strings.TrimSpace(req.Color),
	})
	if err != nil {
		return ShiftType{}, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateRequest) (ShiftType, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, ErrShiftTypeNotFound) {
			return ShiftType{}, apperror.New(apperror.ErrNotFound, "shift type not found")
		}
		return ShiftType{}, err
	}

	fields := make(map[string]any)
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return ShiftType{}, apperror.New(apperror.ErrBadRequest, "name cannot be empty")
		}
		fields["name"] = name
	}
	if req.StartTime != nil || req.EndTime != nil {
		current, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return ShiftType{}, err
		}
		start := current.StartTime
		end := current.EndTime
		if req.StartTime != nil {
			start = *req.StartTime
		}
		if req.EndTime != nil {
			end = *req.EndTime
		}
		start, end, err = normalizeTimes(start, end)
		if err != nil {
			return ShiftType{}, err
		}
		fields["start_time"] = start
		fields["end_time"] = end
	}
	if req.Color != nil {
		fields["color"] = strings.TrimSpace(*req.Color)
	}

	if err := s.repo.Update(ctx, id, fields); err != nil {
		return ShiftType{}, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, ErrShiftTypeNotFound) {
			return apperror.New(apperror.ErrNotFound, "shift type not found")
		}
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func normalizeTimes(startRaw, endRaw string) (string, string, error) {
	start, err := time.Parse("15:04", startRaw)
	if err != nil {
		return "", "", apperror.New(apperror.ErrBadRequest, "startTime must be in HH:MM format")
	}
	end, err := time.Parse("15:04", endRaw)
	if err != nil {
		return "", "", apperror.New(apperror.ErrBadRequest, "endTime must be in HH:MM format")
	}
	if !end.After(start) {
		return "", "", apperror.New(apperror.ErrBadRequest, "endTime must be after startTime")
	}
	return start.Format("15:04:00"), end.Format("15:04:00"), nil
}
