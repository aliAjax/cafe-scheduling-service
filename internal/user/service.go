package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"cafe-scheduling-api/pkg/apperror"
	"cafe-scheduling-api/pkg/auth"
)

type Service struct {
	repo       *Repository
	jwtManager *auth.Manager
}

func NewService(repo *Repository, jwtManager *auth.Manager) *Service {
	return &Service{repo: repo, jwtManager: jwtManager}
}

func (s *Service) Login(ctx context.Context, username, password string) (LoginResponse, error) {
	u, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return LoginResponse{}, apperror.New(apperror.ErrUnauthorized, "invalid username or password")
		}
		return LoginResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return LoginResponse{}, apperror.New(apperror.ErrUnauthorized, "invalid username or password")
	}

	token, err := s.jwtManager.GenerateToken(auth.Claims{
		UserID:     u.ID,
		Username:   u.Username,
		Role:       u.Role,
		EmployeeID: u.EmployeeID,
	})
	if err != nil {
		return LoginResponse{}, fmt.Errorf("generate token: %w", err)
	}

	return LoginResponse{Token: token, User: u}, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (User, error) {
	return s.repo.GetByID(ctx, id)
}
