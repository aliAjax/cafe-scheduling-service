package schedule

import (
	"context"
	"errors"

	"cafe-scheduling-api/internal/businesshours"
	"cafe-scheduling-api/internal/employee"
	"cafe-scheduling-api/internal/shift"
	"cafe-scheduling-api/pkg/apperror"
	"cafe-scheduling-api/pkg/timeutil"
)

type Service struct {
	repo              *Repository
	employeeRepo      *employee.Repository
	shiftRepo         *shift.Repository
	businessHoursRepo *businesshours.Repository
}

func NewService(
	repo *Repository,
	employeeRepo *employee.Repository,
	shiftRepo *shift.Repository,
	businessHoursRepo *businesshours.Repository,
) *Service {
	return &Service{
		repo:              repo,
		employeeRepo:      employeeRepo,
		shiftRepo:         shiftRepo,
		businessHoursRepo: businessHoursRepo,
	}
}

func (s *Service) List(ctx context.Context, weekStart string, employeeID *int64) ([]Schedule, error) {
	start, err := timeutil.ParseDate(weekStart)
	if err != nil {
		return nil, apperror.New(apperror.ErrBadRequest, err.Error())
	}
	end := timeutil.EndOfWeek(start)
	return s.repo.ListByWeek(ctx, timeutil.FormatDate(start), timeutil.FormatDate(end), employeeID)
}

func (s *Service) Get(ctx context.Context, id int64) (Schedule, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrScheduleNotFound) {
			return Schedule{}, apperror.New(apperror.ErrNotFound, "schedule not found")
		}
		return Schedule{}, err
	}
	return item, nil
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (Schedule, error) {
	prepared, err := s.prepare(ctx, req, 0)
	if err != nil {
		return Schedule{}, err
	}
	id, err := s.repo.Create(ctx, prepared)
	if err != nil {
		return Schedule{}, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, req CreateRequest) (Schedule, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, ErrScheduleNotFound) {
			return Schedule{}, apperror.New(apperror.ErrNotFound, "schedule not found")
		}
		return Schedule{}, err
	}
	prepared, err := s.prepare(ctx, req, id)
	if err != nil {
		return Schedule{}, err
	}
	if err := s.repo.Update(ctx, id, prepared); err != nil {
		return Schedule{}, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, ErrScheduleNotFound) {
			return apperror.New(apperror.ErrNotFound, "schedule not found")
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) prepare(ctx context.Context, req CreateRequest, excludeID int64) (Schedule, error) {
	if req.EmployeeID <= 0 || req.ShiftTypeID <= 0 {
		return Schedule{}, apperror.New(apperror.ErrBadRequest, "employeeId and shiftTypeId are required")
	}
	workDate, err := timeutil.ParseDate(req.WorkDate)
	if err != nil {
		return Schedule{}, apperror.New(apperror.ErrBadRequest, err.Error())
	}

	emp, err := s.employeeRepo.GetByID(ctx, req.EmployeeID)
	if err != nil {
		if errors.Is(err, employee.ErrEmployeeNotFound) {
			return Schedule{}, apperror.New(apperror.ErrNotFound, "employee not found")
		}
		return Schedule{}, err
	}
	if !emp.Active {
		return Schedule{}, apperror.New(apperror.ErrConflict, "employee is inactive")
	}

	shiftType, err := s.shiftRepo.GetByID(ctx, req.ShiftTypeID)
	if err != nil {
		if errors.Is(err, shift.ErrShiftTypeNotFound) {
			return Schedule{}, apperror.New(apperror.ErrNotFound, "shift type not found")
		}
		return Schedule{}, err
	}

	day := timeutil.DayOfWeekMondayOne(workDate)
	hours, err := s.businessHoursRepo.GetByDay(ctx, day)
	if err != nil {
		if errors.Is(err, businesshours.ErrBusinessHoursNotFound) {
			return Schedule{}, apperror.New(apperror.ErrConflict, "business hours are not configured for this day")
		}
		return Schedule{}, err
	}
	if hours.IsClosed {
		return Schedule{}, apperror.New(apperror.ErrConflict, "the store is closed on this day")
	}
	if !withinBusinessHours(shiftType.StartTime, shiftType.EndTime, hours.OpenTime, hours.CloseTime) {
		return Schedule{}, apperror.New(apperror.ErrConflict, "shift must be within the business hours for this day")
	}

	overlap, err := s.repo.FindOverlap(ctx, req.EmployeeID, req.WorkDate, shiftType.StartTime, shiftType.EndTime, excludeID)
	if err != nil {
		return Schedule{}, err
	}
	if overlap.ID != 0 {
		return Schedule{}, apperror.Newf(
			apperror.ErrConflict,
			"employee already has an overlapping shift %s-%s",
			overlap.StartTime,
			overlap.EndTime,
		)
	}

	return Schedule{
		EmployeeID:    req.EmployeeID,
		ShiftTypeID:   req.ShiftTypeID,
		WorkDate:      req.WorkDate,
		StartTime:     shiftType.StartTime,
		EndTime:       shiftType.EndTime,
		Hours:         durationHours(shiftType.StartTime, shiftType.EndTime),
		EmployeeName:  emp.Name,
		ShiftTypeName: shiftType.Name,
	}, nil
}

func withinBusinessHours(start, end, open, close string) bool {
	startMinutes := timeToMinutes(start)
	endMinutes := timeToMinutes(end)
	openMinutes := timeToMinutes(open)
	closeMinutes := timeToMinutes(close)
	return startMinutes >= openMinutes && endMinutes <= closeMinutes
}
