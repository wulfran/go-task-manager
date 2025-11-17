package services

import (
	"context"
	"fmt"
	"task-manager/internal/helpers"
	"task-manager/internal/models"
	"task-manager/internal/repository"
)

type UserService interface {
	RegisterUser(ctx context.Context, p models.CreateUserPayload) error
}

type userService struct {
	r repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{
		r: r,
	}
}

func (s userService) RegisterUser(ctx context.Context, p models.CreateUserPayload) error {
	emailExists, err := s.r.CheckIfEmailExists(p.Email)
	if err != nil {
		return fmt.Errorf("RegisterUser: failed to check if email is unique: %v", err)
	}
	if emailExists {
		return fmt.Errorf("RegisterUser: email already in use")
	}
	p.Password = helpers.HashPassword(p.Password)

	if err := s.r.CreateUser(ctx, p); err != nil {
		return fmt.Errorf("RegisterUser: failed to run create user: %v", err)
	}

	return nil
}
