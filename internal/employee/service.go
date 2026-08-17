package employee

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"cafe-scheduling-api/internal/user"
	"cafe-scheduling-api/pkg/apperror"
)

type Service struct {
	repo     *Repository
	userRepo *user.Repository
}

func NewService(repo *Repository, userRepo *user.Repository) *Service {
	return &Service{repo: repo, userRepo: userRepo}
}

func (s *Service) List(ctx context.Context) ([]Employee, error) {
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id int64) (Employee, error) {
	e, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEmployeeNotFound) {
			return Employee{}, apperror.New(apperror.ErrNotFound, "employee not found")
		}
		return Employee{}, err
	}
	return e, nil
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (Employee, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Username = strings.TrimSpace(req.Username)
	if req.Name == "" {
		return Employee{}, apperror.New(apperror.ErrBadRequest, "name is required")
	}
	if (req.Username == "") != (req.Password == "") {
		return Employee{}, apperror.New(apperror.ErrBadRequest, "username and password must be provided together")
	}
	if req.Username != "" && len(req.Password) < 6 {
		return Employee{}, apperror.New(apperror.ErrBadRequest, "password must be at least 6 characters")
	}

	id, err := s.repo.Create(ctx, Employee{Name: req.Name, Phone: req.Phone, Position: req.Position})
	if err != nil {
		return Employee{}, err
	}

	if req.Username != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return Employee{}, fmt.Errorf("hash password: %w", err)
		}
		_, err = s.userRepo.Create(ctx, user.User{
			Username:     req.Username,
			PasswordHash: string(hash),
			Role:         "employee",
			EmployeeID:   &id,
		})
		if err != nil {
			return Employee{}, apperror.New(apperror.ErrConflict, "employee created but user account creation failed: "+err.Error())
		}
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateRequest) (Employee, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, ErrEmployeeNotFound) {
			return Employee{}, apperror.New(apperror.ErrNotFound, "employee not found")
		}
		return Employee{}, err
	}

	fields := make(map[string]any)
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return Employee{}, apperror.New(apperror.ErrBadRequest, "name cannot be empty")
		}
		fields["name"] = name
	}
	if req.Phone != nil {
		fields["phone"] = strings.TrimSpace(*req.Phone)
	}
	if req.Position != nil {
		fields["position"] = strings.TrimSpace(*req.Position)
	}
	if req.Active != nil {
		fields["active"] = *req.Active
	}

	if err := s.repo.Update(ctx, id, fields); err != nil {
		return Employee{}, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, ErrEmployeeNotFound) {
			return apperror.New(apperror.ErrNotFound, "employee not found")
		}
		return err
	}
	return s.repo.SoftDelete(ctx, id)
}
