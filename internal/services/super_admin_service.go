package services

import (
	"context"
	"cws-backend/internal/models"
	"cws-backend/internal/repository"
	"errors"
)

type SuperAdminService struct {
	repo repository.SuperAdminRepo
}

func NewSuperAdminService(repo repository.SuperAdminRepo) *SuperAdminService {
	return &SuperAdminService{
		repo: repo,
	}
}

func (s *SuperAdminService) validateCredentials(ctx context.Context, email string, password string) (models.SuperAdminUser, error) {
	user, err := s.repo.GetByEmail(ctx, email, password)
	if err != nil {
		return models.SuperAdminUser{}, err
	}
	if user.Email == "" {
		return models.SuperAdminUser{}, errors.New("user not found")
	}
	return user, nil
}

func (s *SuperAdminService) GetByEmail(ctx context.Context, email string, password string) (models.SuperAdminUser, error) {
	return s.validateCredentials(ctx, email, password)
}
