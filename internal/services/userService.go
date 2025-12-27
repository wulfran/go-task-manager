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
	LoginUser(p models.LoginPayload) (models.User, error)
	CheckIfEmailExists(e string) (bool, error)
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
	p.Password, err = helpers.HashPassword(p.Password)
	if err != nil {
		return fmt.Errorf("RegisterUser: failed to hash a password, %s", err)
	}

	if err := s.r.CreateUser(ctx, p); err != nil {
		return fmt.Errorf("RegisterUser: failed to run create user: %v", err)
	}

	return nil
}

func (s userService) LoginUser(p models.LoginPayload) (models.User, error) {
	u, err := s.r.GetUserData(p)
	if err != nil {
		return models.User{}, fmt.Errorf("LoginUser: failet to get user data: %v", err)
	}

	return u, nil
}

func (s userService) CheckIfEmailExists(e string) (bool, error) {
	return s.r.CheckIfEmailExists(e)
}
