package services

import (
	"fmt"
	"task-manager/internal/env"
	"task-manager/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	CreateToken(user models.User) (string, error)
}

type authService struct {
	jwtSecret []byte
}

func NewAuthService() AuthService {
	return &authService{
		jwtSecret: []byte(env.Get("jwt_secret")),
	}
}

func (a authService) CreateToken(u models.User) (string, error) {
	claims := jwt.MapClaims{
		"username": u.Email,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("CreateToken: failed to create token string: %v", err)
	}

	return stringToken, nil
}
